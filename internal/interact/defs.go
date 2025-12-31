package interact

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/fuzzy"
	"github.com/dismint/dispass/internal/passio"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
	"github.com/google/uuid"
)

type SearchKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}
type NavKeyMap struct {
	Quit   key.Binding
	Search key.Binding
	Clear  key.Binding
	Nav    key.Binding
	Copy   key.Binding
	Edit   key.Binding
	New    key.Binding
	Del    key.Binding
}
type ViewportKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Save key.Binding
	Next key.Binding
	Prev key.Binding
}

func (k SearchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Confirm}
}
func (k NavKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Search, k.Clear, k.Nav, k.Copy, k.Edit, k.New, k.Del}
}
func (k ViewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Back, k.Save, k.Next, k.Prev}
}

func (k SearchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Confirm},
	}
}
func (k NavKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Search, k.Clear, k.Nav},
		{k.Copy, k.Edit, k.New, k.Del},
	}
}
func (k ViewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Back, k.Save},
		{k.Next, k.Prev},
	}
}

var searchKeyMap = SearchKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm"),
	),
}
var navKeyMap = NavKeyMap{
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
var viewportKeyMap = ViewportKeyMap{
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

type Model struct {
	keyMap    help.KeyMap
	helpModel help.Model

	mode Mode

	viewportSourceInput   textinput.Model
	viewportUsernameInput textinput.Model
	viewportPasswordInput textinput.Model
	viewportUUID          string

	lastQuery       string
	keyInput        textinput.Model
	resultPaginator paginator.Model
	resultLocOnPage int
	topIDs          []string
}

func Initial() Model {
	keyInput := uconst.NewTextInput("search/")
	keyInput.ShowSuggestions = true

	viewportSourceInput := uconst.NewTextInput("Source    ")
	viewportUsernameInput := uconst.NewTextInput("Username  ")
	viewportPasswordInput := uconst.NewTextInput("Password  ")
	viewportPasswordInput.EchoMode = textinput.EchoPassword
	viewportPasswordInput.EchoCharacter = uconst.PasswordChar

	resultPaginator := paginator.New()
	resultPaginator.Type = paginator.Dots
	resultPaginator.PerPage = 10
	resultPaginator.ActiveDot = uconst.SymbolStyle.Render(uconst.PaginatorDotString)
	resultPaginator.InactiveDot = uconst.TextStyle.Render(uconst.PaginatorDotString)

	helpModel := help.New()
	helpModel.Styles = uconst.HelpStyles
	helpModel.ShowAll = true

	return Model{
		keyMap:    navKeyMap,
		helpModel: helpModel,

		mode: ModeNav,

		viewportSourceInput:   viewportSourceInput,
		viewportUsernameInput: viewportUsernameInput,
		viewportPasswordInput: viewportPasswordInput,
		// viewportUUID

		// lastQuery
		keyInput:        keyInput,
		resultPaginator: resultPaginator,
		// resultLocOnPage
		topIDs: make([]string, 0),
	}
}

func (m *Model) populateSuggestions(sm *state.Model) {
	suggestions := make([]string, 0)
	for _, ci := range sm.KeyToCredInfo {
		suggestions = append(suggestions, ci.Source)
		suggestions = append(suggestions, ci.Username)
	}
	m.keyInput.SetSuggestions(suggestions)
}

func (m *Model) populateTopIDs(sm *state.Model, force bool) {
	query := strings.TrimSpace(m.keyInput.Value())
	lastQuery := m.lastQuery
	if query != lastQuery ||
		len(m.topIDs) == 0 ||
		force {
		topIDs := fuzzy.QueryTopIDs(sm, query)
		m.topIDs = topIDs
		if len(topIDs) == 0 {
			m.resultPaginator.TotalPages = 1
		} else {
			m.resultPaginator.SetTotalPages(len(topIDs))
		}
		m.resultPaginator.Page = 0
		m.lastQuery = query
		m.resultLocOnPage = 0
	}

	if credInfo, _, exists := m.getSelectedCredInfo(sm); exists {
		if m.mode == ModeViewport {
			return
		}
		m.viewportSourceInput.SetValue(credInfo.Source)
		m.viewportUsernameInput.SetValue(credInfo.Username)
		m.viewportPasswordInput.SetValue(credInfo.Password)
		m.viewportSourceInput.CursorEnd()
		m.viewportUsernameInput.CursorEnd()
		m.viewportPasswordInput.CursorEnd()
	}
}

func (m *Model) getSelectedCredInfo(sm *state.Model) (state.CredInfo, string, bool) {
	if len(m.topIDs) == 0 {
		return state.CredInfo{}, "", false
	}

	start, _ := m.resultPaginator.GetSliceBounds(len(m.topIDs))
	id := m.topIDs[start+m.resultLocOnPage]
	credInfo, exists := sm.KeyToCredInfo[id]
	return credInfo, id, exists
}

func (m *Model) viewViewport(sm *state.Model) string {
	_, _, exists := m.getSelectedCredInfo(sm)
	if !exists && m.mode != ModeViewport {
		return "Feeling empty, create new credentials?"
	}
	return fmt.Sprintf("%v\n%v\n%v",
		m.viewportSourceInput.View(),
		m.viewportUsernameInput.View(),
		m.viewportPasswordInput.View(),
	)
}

func (m *Model) Update(msg tea.Msg, sm *state.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	var cmd tea.Cmd
	m.keyInput, cmd = m.keyInput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewportSourceInput, cmd = m.viewportSourceInput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewportUsernameInput, cmd = m.viewportUsernameInput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewportPasswordInput, cmd = m.viewportPasswordInput.Update(msg)
	cmds = append(cmds, cmd)
	// normally would update pagination, but we want to override keybinds

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, searchKeyMap.Quit):
			sm.Quitting = true
			cmds = append(cmds, tea.Quit)
		// search mode
		case m.mode == ModeSearch && key.Matches(typedMsg, searchKeyMap.Confirm):
			m.keyInput.Blur()
			m.mode = ModeNav
			m.keyMap = navKeyMap
			m.helpModel.ShowAll = true
		// nav mode
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Search):
			cmds = append(cmds, m.keyInput.Focus())
			m.mode = ModeSearch
			m.keyMap = searchKeyMap
			m.helpModel.ShowAll = false
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Clear):
			m.keyInput.SetValue("")
			m.keyInput.CursorEnd()
			m.populateSuggestions(sm)
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Nav):
			direction := typedMsg.String()
			switch direction {
			case "left", "h":
				m.resultPaginator.PrevPage()
				m.resultLocOnPage = 0
			case "right", "l":
				m.resultPaginator.NextPage()
				m.resultLocOnPage = 0
			case "up", "k":
				m.resultLocOnPage = max(
					m.resultLocOnPage-1,
					0,
				)
			default:
				start, end := m.resultPaginator.GetSliceBounds(len(m.topIDs))
				m.resultLocOnPage = min(
					m.resultLocOnPage+1,
					end-start-1,
				)
			}
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Copy):
			if credInfo, _, exists := m.getSelectedCredInfo(sm); exists {
				clipboard.WriteAll(credInfo.Password)
				cmds = append(cmds, state.NotificationMsg(
					"Credentials Copied",
					state.MessageLevelSuccess,
				))
			}
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Edit):
			if _, _, exists := m.getSelectedCredInfo(sm); exists {
				cmds = append(cmds, m.viewportSourceInput.Focus())
				m.mode = ModeViewport
				m.keyMap = viewportKeyMap
			}
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.New):
			cmds = append(cmds, m.viewportSourceInput.Focus())
			m.mode = ModeViewport
			m.keyMap = viewportKeyMap
			m.viewportUUID = uuid.NewString()
			m.viewportSourceInput.SetValue("")
			m.viewportUsernameInput.SetValue("")
			m.viewportPasswordInput.SetValue("")
			m.viewportSourceInput.CursorEnd()
			m.viewportUsernameInput.CursorEnd()
			m.viewportPasswordInput.CursorEnd()
		case m.mode == ModeNav && key.Matches(typedMsg, navKeyMap.Del):
			if _, id, exists := m.getSelectedCredInfo(sm); exists {
				fuzzy.RemoveFuzzy(sm, id)
				delete(sm.KeyToCredInfo, id)
				passio.WriteStateCreds(sm)
				m.populateTopIDs(sm, true)
				m.populateSuggestions(sm)
				cmds = append(cmds, state.NotificationMsg(
					"Credentials Deleted",
					state.MessageLevelSuccess,
				))
			}
		// viewport mode
		case m.mode == ModeViewport && key.Matches(typedMsg, viewportKeyMap.Back):
			m.viewportSourceInput.Blur()
			m.viewportUsernameInput.Blur()
			m.viewportPasswordInput.Blur()
			m.viewportSourceInput.SetValue("")
			m.viewportUsernameInput.SetValue("")
			m.viewportPasswordInput.SetValue("")
			m.viewportSourceInput.CursorEnd()
			m.viewportUsernameInput.CursorEnd()
			m.viewportPasswordInput.CursorEnd()
			m.viewportUUID = ""
			m.mode = ModeNav
			m.keyMap = navKeyMap
		case m.mode == ModeViewport && key.Matches(typedMsg, viewportKeyMap.Prev):
			if m.viewportUsernameInput.Focused() {
				m.viewportUsernameInput.Blur()
				cmds = append(cmds, m.viewportSourceInput.Focus())
			} else if m.viewportPasswordInput.Focused() {
				m.viewportPasswordInput.Blur()
				cmds = append(cmds, m.viewportUsernameInput.Focus())
			} else if m.viewportSourceInput.Focused() {
				m.viewportSourceInput.Blur()
				cmds = append(cmds, m.viewportPasswordInput.Focus())
			}
			m.viewportSourceInput.CursorEnd()
			m.viewportUsernameInput.CursorEnd()
			m.viewportPasswordInput.CursorEnd()
		case m.mode == ModeViewport && key.Matches(typedMsg, viewportKeyMap.Next):
			if m.viewportUsernameInput.Focused() {
				m.viewportUsernameInput.Blur()
				cmds = append(cmds, m.viewportPasswordInput.Focus())
			} else if m.viewportSourceInput.Focused() {
				m.viewportSourceInput.Blur()
				cmds = append(cmds, m.viewportUsernameInput.Focus())
			} else if m.viewportPasswordInput.Focused() {
				m.viewportPasswordInput.Blur()
				cmds = append(cmds, m.viewportSourceInput.Focus())
			}
			m.viewportSourceInput.CursorEnd()
			m.viewportUsernameInput.CursorEnd()
			m.viewportPasswordInput.CursorEnd()
		case m.mode == ModeViewport && key.Matches(typedMsg, viewportKeyMap.Save):
			id := m.viewportUUID
			if id == "" {
				if _, existingId, exists := m.getSelectedCredInfo(sm); exists {
					id = existingId
				} else {
					log.Fatalf("no existing selection when one needed")
				}
			}
			ci := state.CredInfo{
				Source:   m.viewportSourceInput.Value(),
				Username: m.viewportUsernameInput.Value(),
				Password: m.viewportPasswordInput.Value(),
			}
			fuzzy.UpdateFuzzy(sm, id, ci)
			sm.KeyToCredInfo[id] = ci
			passio.WriteStateCreds(sm)
			m.populateTopIDs(sm, true)
			m.populateSuggestions(sm)
			// also do a back...
			m.viewportSourceInput.Blur()
			m.viewportUsernameInput.Blur()
			m.viewportPasswordInput.Blur()
			m.viewportSourceInput.SetValue("")
			m.viewportUsernameInput.SetValue("")
			m.viewportPasswordInput.SetValue("")
			m.viewportSourceInput.CursorEnd()
			m.viewportUsernameInput.CursorEnd()
			m.viewportPasswordInput.CursorEnd()
			m.viewportUUID = ""
			m.mode = ModeNav
			m.keyMap = navKeyMap

			cmds = append(cmds, state.NotificationMsg(
				"Credentials Saved",
				state.MessageLevelSuccess,
			))
		}
	}

	if _, ok := msg.(tea.KeyMsg); ok {
		m.populateTopIDs(sm, false)
	}
	if sm.Dirty {
		m.populateTopIDs(sm, true)
		m.populateSuggestions(sm)
		sm.Dirty = false
	}

	return tea.Batch(cmds...)
}

func (m *Model) View(sm *state.Model) string {
	start, end := m.resultPaginator.GetSliceBounds(len(m.topIDs))
	topPageIDs := m.topIDs[start:end]

	var resultList string
	for locOnPage, topID := range topPageIDs {
		prefix := " "
		if locOnPage == m.resultLocOnPage {
			prefix = uconst.SymbolStyle.Render(">")
		}
		resultList += fmt.Sprintf("%v %v %v\n",
			prefix,
			uconst.TruncAndPadListElem(sm.KeyToCredInfo[topID].Source),
			uconst.TruncAndPadListElem(sm.KeyToCredInfo[topID].Username),
		)
	}
	if resultList == "" {
		resultList = "No Results Found"
	}

	view := fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v\n\n%v",
		m.helpModel.View(m.keyMap),
		m.keyInput.View(),
		uconst.ViewportViewStyle.Render(m.viewViewport(sm)),
		m.resultPaginator.View(),
		resultList,
	)

	return uconst.ViewStyle.Render(view)
}
