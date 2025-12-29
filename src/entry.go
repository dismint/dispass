package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type entryKeyMap struct {
	Quit  key.Binding
	Enter key.Binding
}

func (k entryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Enter}
}

func (k entryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Enter},
	}
}

var entryKeys = entryKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "enter"),
	),
}

type EntryModel struct {
	keys entryKeyMap
	help help.Model

	confirming bool
	completed  bool

	passwordInput        textinput.Model
	confirmPasswordInput textinput.Model
}

func initialEntryModel(m *model) {
	passwordInput := textinput.New()
	passwordInput.Focus()
	passwordInput.CharLimit = -1
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '▪'
	passwordInput.Prompt = "$ "
	passwordInput.Cursor.Style = symbolStyle
	passwordInput.PromptStyle = symbolStyle
	passwordInput.TextStyle = textStyle

	confirmPasswordInput := textinput.New()
	confirmPasswordInput.CharLimit = -1
	confirmPasswordInput.EchoMode = textinput.EchoPassword
	confirmPasswordInput.EchoCharacter = '▪'
	confirmPasswordInput.Prompt = "$ "
	confirmPasswordInput.Cursor.Style = symbolStyle
	confirmPasswordInput.PromptStyle = symbolStyle
	confirmPasswordInput.TextStyle = textStyle

	entryHelp := help.New()
	entryHelp.Styles = helpStyles

	m.entryModel = EntryModel{
		keys: entryKeys,
		help: entryHelp,

		confirming: false,
		completed:  false,

		passwordInput:        passwordInput,
		confirmPasswordInput: confirmPasswordInput,
	}
}

func (m *model) populateSuggestions() {
	suggestions := make([]string, 0)
	for _, credInfo := range m.keyToCredInfo {
		suggestions = append(suggestions, credInfo.Source)
		suggestions = append(suggestions, credInfo.Username)
	}
	m.selectModel.keyInput.SetSuggestions(suggestions)
}

func (m *model) passwordComplete(createNew bool) (tea.Cmd, error) {
	cmds := make([]tea.Cmd, 0)
	m.keyToCredInfo = make(map[string]credInfo)
	m.secret = secretFromString(m.entryModel.passwordInput.Value())
	if createNew {
		m.writePass()
	}
	if err := m.readPass(); err != nil {
		// this error should only happen when we give the wrong password and can't decrypt
		return showNotification("Incorrect Password", MessageLevelError), err
	}
	m.screen = SelectScreen
	m.entryModel.completed = true

	// set up selectModel
	m.populateSuggestions()
	m.initFuzzy()
	m.populateTopIDs(true)

	// turn off
	m.entryModel.passwordInput.Blur()
	m.entryModel.confirmPasswordInput.Blur()

	return tea.Batch(cmds...), nil
}

func updateEntryModel(m *model, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, entryKeys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, entryKeys.Enter):
			if m.entryModel.confirming {
				if m.entryModel.passwordInput.Value() != m.entryModel.confirmPasswordInput.Value() {
					break
				}
				pcCmds, err := m.passwordComplete(true)
				if err != nil {
					cmds = append(cmds, pcCmds)
					break
				}
				cmds = append(cmds, pcCmds)
			} else {
				if _, err := os.Stat(dataFileName); err == nil {
					pcCmds, err := m.passwordComplete(false)
					if err != nil {
						cmds = append(cmds, pcCmds)
						break
					}
					cmds = append(cmds, pcCmds)
				} else if os.IsNotExist(err) {
					m.entryModel.confirming = true
					cmds = append(cmds, m.entryModel.confirmPasswordInput.Focus())
					m.entryModel.passwordInput.Blur()
				} else {
					log.Errorf("error reading file: %v", err)
					panic(err)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.entryModel.passwordInput, cmd = m.entryModel.passwordInput.Update(msg)
	cmds = append(cmds, cmd)
	m.entryModel.confirmPasswordInput, cmd = m.entryModel.confirmPasswordInput.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func viewEntryModel(m *model) string {
	view := fmt.Sprintf("%v\n\n",
		m.entryModel.help.View(m.entryModel.keys),
	)
	view += fmt.Sprintf("%v",
		m.entryModel.passwordInput.View(),
	)
	if m.entryModel.confirming {
		view += fmt.Sprintf("\n%v",
			m.entryModel.confirmPasswordInput.View(),
		)
	}
	return finalWrapStyle.Render(view)
}
