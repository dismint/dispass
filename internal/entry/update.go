package entry

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/fuzzy"
	"github.com/dismint/dispass/internal/passio"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

func (m *Model) passwordComplete(createNew bool, sm *state.Model) (tea.Cmd, error) {
	sm.Secret = passio.SecretFromString(m.passwordInput.Value())
	if createNew {
		passio.WriteStateCreds(sm)
	}
	if err := passio.ReadStateCreds(sm); err != nil {
		// this error should only happen when we give the wrong password and can't decrypt
		return state.NotificationMsg("Incorrect Password", state.MessageLevelError), err
	}
	sm.Screen = state.InteractScreen
	sm.Dirty = true

	fuzzy.InitFuzzy(sm)

	m.passwordInput.Blur()
	m.confirmPasswordInput.Blur()

	return nil, nil
}

func (m *Model) Update(msg tea.Msg, sm *state.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	var cmd tea.Cmd
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	cmds = append(cmds, cmd)
	m.confirmPasswordInput, cmd = m.confirmPasswordInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Quit):
			sm.Quitting = true
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, keyMap.Enter):
			if m.confirming {
				// data does not exist and we confirmed password, try decrypting
				if m.passwordInput.Value() != m.confirmPasswordInput.Value() {
					break
				}
				pcCmd, err := m.passwordComplete(true, sm)
				if err != nil {
					cmds = append(cmds, pcCmd)
					break
				}
			} else {
				// first entry, check which scenario we're in
				if _, err := os.Stat(uconst.DataFileName); err == nil {
					// data exists, try decrypting
					pcCmd, err := m.passwordComplete(false, sm)
					if err != nil {
						cmds = append(cmds, pcCmd)
						break
					}
				} else if os.IsNotExist(err) {
					// data does not exist, confirm password
					m.confirming = true
					cmds = append(cmds, m.confirmPasswordInput.Focus())
					m.passwordInput.Blur()
				} else {
					log.Fatalf("error reading file: %v", err)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}
