package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"golang.design/x/clipboard"
)

type searchKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}
type navKeyMap struct {
	Quit   key.Binding
	Search key.Binding
	Clear  key.Binding
	Nav    key.Binding
	Copy   key.Binding
	Edit   key.Binding
	New    key.Binding
	Del    key.Binding
}
type viewportKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Save key.Binding
	Next key.Binding
	Prev key.Binding
}

func (k searchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Confirm}
}
func (k navKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Search, k.Clear, k.Nav, k.Copy, k.Edit, k.New, k.Del}
}
func (k viewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Back, k.Save, k.Next, k.Prev}
}

func (k searchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Confirm},
	}
}
func (k navKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Search, k.Clear, k.Nav},
		{k.Copy, k.Edit, k.New, k.Del},
	}
}
func (k viewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Back, k.Save},
		{k.Next, k.Prev},
	}
}

var searchKeys = searchKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm"),
	),
}
var navKeys = navKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Search: key.NewBinding(
		key.WithKeys("s", "/"),
		key.WithHelp("s /", "search"),
	),
	Clear: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear"),
	),
	Nav: key.NewBinding(
		key.WithKeys("left", "up", "right", "down", "h", "k", "l", "j"),
		key.WithHelp("←↑→↓", "nav"),
	),
	Copy: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "copy"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Del: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}
var viewportKeys = viewportKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "save"),
	),
	Next: key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("↓ tab", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "prev"),
	),
}

type Mode int

const (
	ModeSearch Mode = iota
	ModeNav
	ModeViewport
)

type SelectModel struct {
	keys help.KeyMap
	help help.Model

	lastQuery string
	mode      Mode

	viewportSourceInput   textinput.Model
	viewportUsernameInput textinput.Model
	viewportPasswordInput textinput.Model
	viewportUUID          string

	keyInput        textinput.Model
	resultPaginator paginator.Model
	resultLocOnPage int
	topIDs          []string
}

func initialSelectModel(m *model) {
	keyInput := newTextInput("search/")
	keyInput.ShowSuggestions = true

	viewportSourceInput := newTextInput("Source    ")
	viewportUsernameInput := newTextInput("Username  ")
	viewportPasswordInput := newTextInput("Password  ")
	viewportPasswordInput.EchoMode = textinput.EchoPassword
	viewportPasswordInput.EchoCharacter = '▪'

	resultPaginator := paginator.New()
	resultPaginator.Type = paginator.Dots
	resultPaginator.PerPage = 10
	resultPaginator.ActiveDot = symbolStyle.Render("▪")
	resultPaginator.InactiveDot = textStyle.Render("▪")

	selectHelp := help.New()
	selectHelp.Styles = helpStyles
	selectHelp.ShowAll = true

	m.selectModel = SelectModel{
		keys: navKeys,
		help: selectHelp,

		mode: ModeNav,

		viewportSourceInput:   viewportSourceInput,
		viewportUsernameInput: viewportUsernameInput,
		viewportPasswordInput: viewportPasswordInput,

		keyInput:        keyInput,
		resultPaginator: resultPaginator,
	}
}

func (m *model) populateTopIDs(manualAlways bool) {
	query := strings.TrimSpace(m.selectModel.keyInput.Value())
	lastQuery := m.selectModel.lastQuery
	if query != lastQuery ||
		len(m.selectModel.topIDs) == 0 ||
		manualAlways {
		topIDs := queryTopIDs(query)
		m.selectModel.topIDs = topIDs
		if len(topIDs) == 0 {
			m.selectModel.resultPaginator.TotalPages = 1
		} else {
			m.selectModel.resultPaginator.SetTotalPages(len(topIDs))
		}
		m.selectModel.resultPaginator.Page = 0
		m.selectModel.lastQuery = query
		m.selectModel.resultLocOnPage = 0
	}

	if credInfo, _, exists := m.getSelectedCredInfo(); exists {
		if m.selectModel.mode == ModeViewport {
			return
		}
		m.selectModel.viewportSourceInput.SetValue(credInfo.Source)
		m.selectModel.viewportUsernameInput.SetValue(credInfo.Username)
		m.selectModel.viewportPasswordInput.SetValue(credInfo.Password)
		m.selectModel.viewportSourceInput.CursorEnd()
		m.selectModel.viewportUsernameInput.CursorEnd()
		m.selectModel.viewportPasswordInput.CursorEnd()
	}
}

func (m *model) getSelectedCredInfo() (credInfo, string, bool) {
	if len(m.selectModel.topIDs) == 0 {
		return credInfo{}, "", false
	}

	start, _ := m.selectModel.resultPaginator.GetSliceBounds(len(m.selectModel.topIDs))
	id := m.selectModel.topIDs[start+m.selectModel.resultLocOnPage]
	credInfo, exists := m.keyToCredInfo[id]
	return credInfo, id, exists
}

func (m *model) renderViewport() string {
	_, _, exists := m.getSelectedCredInfo()
	if !exists && m.selectModel.mode != ModeViewport {
		return "Feeling empty, create new credentials?"
	}
	return fmt.Sprintf("%v\n%v\n%v",
		m.selectModel.viewportSourceInput.View(),
		m.selectModel.viewportUsernameInput.View(),
		m.selectModel.viewportPasswordInput.View(),
	)
}

func updateSelectModel(m *model, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	var cmd tea.Cmd
	m.selectModel.keyInput, cmd = m.selectModel.keyInput.Update(msg)
	cmds = append(cmds, cmd)
	m.selectModel.viewportSourceInput, cmd = m.selectModel.viewportSourceInput.Update(msg)
	cmds = append(cmds, cmd)
	m.selectModel.viewportUsernameInput, cmd = m.selectModel.viewportUsernameInput.Update(msg)
	cmds = append(cmds, cmd)
	m.selectModel.viewportPasswordInput, cmd = m.selectModel.viewportPasswordInput.Update(msg)
	cmds = append(cmds, cmd)
	// normally would update pagination, but we want to override keybinds

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, searchKeys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		// search mode
		case m.selectModel.mode == ModeSearch && key.Matches(typedMsg, searchKeys.Confirm):
			m.selectModel.keyInput.Blur()
			m.selectModel.mode = ModeNav
			m.selectModel.keys = navKeys
			m.selectModel.help.ShowAll = true
		// nav mode
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Search):
			cmds = append(cmds, m.selectModel.keyInput.Focus())
			m.selectModel.mode = ModeSearch
			m.selectModel.keys = searchKeys
			m.selectModel.help.ShowAll = false
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Clear):
			m.selectModel.keyInput.SetValue("")
			m.selectModel.keyInput.CursorEnd()
			m.populateSuggestions()
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Nav):
			direction := typedMsg.String()
			switch direction {
			case "left", "h":
				m.selectModel.resultPaginator.PrevPage()
				m.selectModel.resultLocOnPage = 0
			case "right", "l":
				m.selectModel.resultPaginator.NextPage()
				m.selectModel.resultLocOnPage = 0
			case "up", "k":
				m.selectModel.resultLocOnPage = max(
					m.selectModel.resultLocOnPage-1,
					0,
				)
			default:
				start, end := m.selectModel.resultPaginator.GetSliceBounds(len(m.selectModel.topIDs))
				m.selectModel.resultLocOnPage = min(
					m.selectModel.resultLocOnPage+1,
					end-start-1,
				)
			}
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Copy):
			if credInfo, _, exists := m.getSelectedCredInfo(); exists {
				clipboard.Write(clipboard.FmtText, []byte(credInfo.Password))
				cmds = append(cmds, showNotification("Credentials Copied", MessageLevelSuccess))
			}
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Edit):
			if _, _, exists := m.getSelectedCredInfo(); exists {
				cmds = append(cmds, m.selectModel.viewportSourceInput.Focus())
				m.selectModel.mode = ModeViewport
				m.selectModel.keys = viewportKeys
			}
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.New):
			cmds = append(cmds, m.selectModel.viewportSourceInput.Focus())
			m.selectModel.mode = ModeViewport
			m.selectModel.keys = viewportKeys
			m.selectModel.viewportUUID = uuid.NewString()
			m.selectModel.viewportSourceInput.SetValue("")
			m.selectModel.viewportUsernameInput.SetValue("")
			m.selectModel.viewportPasswordInput.SetValue("")
			m.selectModel.viewportSourceInput.CursorEnd()
			m.selectModel.viewportUsernameInput.CursorEnd()
			m.selectModel.viewportPasswordInput.CursorEnd()
		case m.selectModel.mode == ModeNav && key.Matches(typedMsg, navKeys.Del):
			if _, id, exists := m.getSelectedCredInfo(); exists {
				m.removeFuzzy(id)
				delete(m.keyToCredInfo, id)
				m.writePass()
				m.populateTopIDs(true)
				m.populateSuggestions()
				cmds = append(cmds, showNotification("Credentials Deleted", MessageLevelSuccess))
			}
		// viewport mode
		case m.selectModel.mode == ModeViewport && key.Matches(typedMsg, viewportKeys.Back):
			m.selectModel.viewportSourceInput.Blur()
			m.selectModel.viewportUsernameInput.Blur()
			m.selectModel.viewportPasswordInput.Blur()
			m.selectModel.viewportSourceInput.SetValue("")
			m.selectModel.viewportUsernameInput.SetValue("")
			m.selectModel.viewportPasswordInput.SetValue("")
			m.selectModel.viewportSourceInput.CursorEnd()
			m.selectModel.viewportUsernameInput.CursorEnd()
			m.selectModel.viewportPasswordInput.CursorEnd()
			m.selectModel.viewportUUID = ""
			m.selectModel.mode = ModeNav
			m.selectModel.keys = navKeys
		case m.selectModel.mode == ModeViewport && key.Matches(typedMsg, viewportKeys.Prev):
			if m.selectModel.viewportUsernameInput.Focused() {
				m.selectModel.viewportUsernameInput.Blur()
				cmds = append(cmds, m.selectModel.viewportSourceInput.Focus())
			} else if m.selectModel.viewportPasswordInput.Focused() {
				m.selectModel.viewportPasswordInput.Blur()
				cmds = append(cmds, m.selectModel.viewportUsernameInput.Focus())
			} else if m.selectModel.viewportSourceInput.Focused() {
				m.selectModel.viewportSourceInput.Blur()
				cmds = append(cmds, m.selectModel.viewportPasswordInput.Focus())
			}
			m.selectModel.viewportSourceInput.CursorEnd()
			m.selectModel.viewportUsernameInput.CursorEnd()
			m.selectModel.viewportPasswordInput.CursorEnd()
		case m.selectModel.mode == ModeViewport && key.Matches(typedMsg, viewportKeys.Next):
			if m.selectModel.viewportUsernameInput.Focused() {
				m.selectModel.viewportUsernameInput.Blur()
				cmds = append(cmds, m.selectModel.viewportPasswordInput.Focus())
			} else if m.selectModel.viewportSourceInput.Focused() {
				m.selectModel.viewportSourceInput.Blur()
				cmds = append(cmds, m.selectModel.viewportUsernameInput.Focus())
			} else if m.selectModel.viewportPasswordInput.Focused() {
				m.selectModel.viewportPasswordInput.Blur()
				cmds = append(cmds, m.selectModel.viewportSourceInput.Focus())
			}
			m.selectModel.viewportSourceInput.CursorEnd()
			m.selectModel.viewportUsernameInput.CursorEnd()
			m.selectModel.viewportPasswordInput.CursorEnd()
		case m.selectModel.mode == ModeViewport && key.Matches(typedMsg, viewportKeys.Save):
			id := m.selectModel.viewportUUID
			if id == "" {
				if _, existingId, exists := m.getSelectedCredInfo(); exists {
					id = existingId
				} else {
					log.Fatalf("no existing selection when one needed")
				}
			}
			ci := credInfo{
				Source:   m.selectModel.viewportSourceInput.Value(),
				Username: m.selectModel.viewportUsernameInput.Value(),
				Password: m.selectModel.viewportPasswordInput.Value(),
			}
			m.updateFuzzy(id, ci)
			m.keyToCredInfo[id] = ci
			m.writePass()
			m.populateTopIDs(true)
			m.populateSuggestions()
			// also do a back...
			m.selectModel.viewportSourceInput.Blur()
			m.selectModel.viewportUsernameInput.Blur()
			m.selectModel.viewportPasswordInput.Blur()
			m.selectModel.viewportSourceInput.SetValue("")
			m.selectModel.viewportUsernameInput.SetValue("")
			m.selectModel.viewportPasswordInput.SetValue("")
			m.selectModel.viewportSourceInput.CursorEnd()
			m.selectModel.viewportUsernameInput.CursorEnd()
			m.selectModel.viewportPasswordInput.CursorEnd()
			m.selectModel.viewportUUID = ""
			m.selectModel.mode = ModeNav
			m.selectModel.keys = navKeys

			cmds = append(cmds, showNotification("Credentials Saved", MessageLevelSuccess))
		}
	}

	if _, ok := msg.(tea.KeyMsg); ok {
		m.populateTopIDs(false)
	}

	return tea.Batch(cmds...)
}

func viewSelectModel(m *model) string {
	start, end := m.selectModel.resultPaginator.GetSliceBounds(len(m.selectModel.topIDs))
	topPageIDs := m.selectModel.topIDs[start:end]

	var resultList string
	for locOnPage, topID := range topPageIDs {
		prefix := " "
		if locOnPage == m.selectModel.resultLocOnPage {
			prefix = symbolStyle.Render(">")
		}
		resultList += fmt.Sprintf("%v %v %v\n",
			prefix,
			truncAndPadListElem(m.keyToCredInfo[topID].Source),
			truncAndPadListElem(m.keyToCredInfo[topID].Username),
		)
	}
	if resultList == "" {
		resultList = "No Results Found"
	}

	view := fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v\n\n%v",
		m.selectModel.help.View(m.selectModel.keys),
		m.selectModel.keyInput.View(),
		viewportStyle.Render(m.renderViewport()),
		m.selectModel.resultPaginator.View(),
		resultList,
	)

	return finalWrapStyleBounded.Render(view)
}
