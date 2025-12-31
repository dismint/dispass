package interact

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/fuzzy"
	"github.com/dismint/dispass/internal/passio"
	"github.com/dismint/dispass/internal/state"
	"github.com/google/uuid"
)

func (m *Model) setViewportCredInfo(credInfo state.CredInfo, blur bool) {
	m.viewportSourceInput.SetValue(credInfo.Source)
	m.viewportUsernameInput.SetValue(credInfo.Username)
	m.viewportPasswordInput.SetValue(credInfo.Password)
	m.viewportSourceInput.CursorEnd()
	m.viewportUsernameInput.CursorEnd()
	m.viewportPasswordInput.CursorEnd()
	if blur {
		m.viewportSourceInput.Blur()
		m.viewportUsernameInput.Blur()
		m.viewportPasswordInput.Blur()
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
	if query != m.lastQuery || len(m.topIDs) == 0 || force {
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
		m.setViewportCredInfo(credInfo, false)
	}
}

func (m *Model) updateSearch(keyMsg tea.KeyMsg) {
	switch {
	case key.Matches(keyMsg, searchKeyMap.Confirm):
		m.keyInput.Blur()
		m.mode = ModeNav
		m.keyMap = navKeyMap
		m.helpModel.ShowAll = true
	}
}

func (m *Model) updateNav(keyMsg tea.KeyMsg, sm *state.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch {
	case key.Matches(keyMsg, navKeyMap.Search):
		cmds = append(cmds, m.keyInput.Focus())
		m.mode = ModeSearch
		m.keyMap = searchKeyMap
		m.helpModel.ShowAll = false
	case key.Matches(keyMsg, navKeyMap.Clear):
		m.keyInput.SetValue("")
		m.keyInput.CursorEnd()
	case key.Matches(keyMsg, navKeyMap.Nav):
		switch keyMsg.String() {
		case "left", "h":
			m.resultPaginator.PrevPage()
			m.resultLocOnPage = 0
		case "right", "l":
			m.resultPaginator.NextPage()
			m.resultLocOnPage = 0
		case "up", "k":
			m.resultLocOnPage = max(m.resultLocOnPage-1, 0)
		case "down", "j":
			start, end := m.resultPaginator.GetSliceBounds(len(m.topIDs))
			m.resultLocOnPage = min(m.resultLocOnPage+1, end-start-1)
		}
	case key.Matches(keyMsg, navKeyMap.Copy):
		if credInfo, _, exists := m.getSelectedCredInfo(sm); exists {
			clipboard.WriteAll(credInfo.Password)
			cmds = append(cmds, state.NotificationMsg(
				"Credentials Copied",
				state.MessageLevelSuccess,
			))
		}
	case key.Matches(keyMsg, navKeyMap.Edit):
		if _, _, exists := m.getSelectedCredInfo(sm); exists {
			cmds = append(cmds, m.viewportSourceInput.Focus())
			m.mode = ModeViewport
			m.keyMap = viewportKeyMap
		}
	case key.Matches(keyMsg, navKeyMap.New):
		cmds = append(cmds, m.viewportSourceInput.Focus())
		m.mode = ModeViewport
		m.keyMap = viewportKeyMap
		m.viewportUUID = uuid.NewString()
		m.setViewportCredInfo(
			state.CredInfo{Source: "", Username: "", Password: ""},
			false,
		)
	case key.Matches(keyMsg, navKeyMap.Del):
		if _, id, exists := m.getSelectedCredInfo(sm); exists {
			fuzzy.RemoveFuzzy(sm, id)
			delete(sm.KeyToCredInfo, id)
			passio.WriteStateCreds(sm)
			cmds = append(cmds, state.NotificationMsg(
				"Credentials Deleted",
				state.MessageLevelSuccess,
			))
		}
	}

	return tea.Batch(cmds...)
}

func (m *Model) updateViewport(keyMsg tea.KeyMsg, sm *state.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch {
	case key.Matches(keyMsg, viewportKeyMap.Back):
		m.setViewportCredInfo(
			state.CredInfo{Source: "", Username: "", Password: ""},
			true,
		)
		m.viewportUUID = ""
		m.mode = ModeNav
		m.keyMap = navKeyMap
	case key.Matches(keyMsg, viewportKeyMap.Prev):
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
	case key.Matches(keyMsg, viewportKeyMap.Next):
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
	case key.Matches(keyMsg, viewportKeyMap.Save):
		id := m.viewportUUID
		if id == "" {
			if _, existingId, exists := m.getSelectedCredInfo(sm); exists {
				id = existingId
			} else {
				log.Fatalf("no existing selection when one needed")
			}
		}
		credInfo := state.CredInfo{
			Source:   m.viewportSourceInput.Value(),
			Username: m.viewportUsernameInput.Value(),
			Password: m.viewportPasswordInput.Value(),
		}
		fuzzy.UpdateFuzzy(sm, id, credInfo)
		sm.KeyToCredInfo[id] = credInfo
		passio.WriteStateCreds(sm)
		m.setViewportCredInfo(
			state.CredInfo{Source: "", Username: "", Password: ""},
			true,
		)
		m.viewportUUID = ""
		m.mode = ModeNav
		m.keyMap = navKeyMap

		cmds = append(cmds, state.NotificationMsg(
			"Credentials Saved",
			state.MessageLevelSuccess,
		))
	}

	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg, sm *state.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, ti := range []*textinput.Model{
		&m.keyInput,
		&m.viewportSourceInput,
		&m.viewportUsernameInput,
		&m.viewportPasswordInput,
	} {
		tiPointer, cmd := ti.Update(msg)
		cmds = append(cmds, cmd)
		*ti = tiPointer
	}
	// manually update the paginator in code later

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, searchKeyMap.Quit):
			sm.Quitting = true
			cmds = append(cmds, tea.Quit)
		case m.mode == ModeSearch:
			m.updateSearch(typedMsg)
		case m.mode == ModeNav:
			cmds = append(cmds, m.updateNav(typedMsg, sm))
		case m.mode == ModeViewport:
			cmds = append(cmds, m.updateViewport(typedMsg, sm))
		}
	}

	if _, ok := msg.(tea.KeyMsg); ok {
		m.populateTopIDs(sm, false)
	}
	if sm.Dirty {
		m.populateTopIDs(sm, true)
		m.populateSuggestions(sm)
	}

	return tea.Batch(cmds...)
}
