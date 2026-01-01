package changemaster

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dismint/dispass/internal/passio"
	"github.com/dismint/dispass/internal/state"
)

func (m *Model) transitionState(sm *state.Model) {
	sm.Screen = state.InteractScreen

	m.passwordInput.SetValue("")
	m.confirmPasswordInput.SetValue("")
	m.passwordInput.Blur()
	m.confirmPasswordInput.Blur()
}

func (m *Model) passwordComplete(sm *state.Model) {
	sm.Secret = passio.SecretFromString(m.passwordInput.Value())
	passio.WriteStateCreds(sm)

	m.transitionState(sm)
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
			if m.confirming && m.confirmPasswordInput.Value() != "" {
				if m.passwordInput.Value() != m.confirmPasswordInput.Value() {
					cmds = append(cmds, state.NotificationMsg(
						"Passwords do not match",
						state.MessageLevelError,
					))
					m.passwordInput.SetValue("")
					m.confirmPasswordInput.SetValue("")
					m.confirmPasswordInput.Blur()
					cmds = append(cmds, m.passwordInput.Focus())
					m.confirming = false
				} else {
					m.passwordComplete(sm)
				}
			} else if m.passwordInput.Value() != "" {
				m.confirming = true
				cmds = append(cmds, m.confirmPasswordInput.Focus())
				m.passwordInput.Blur()
			}
		case key.Matches(msg, keyMap.Back):
			m.transitionState(sm)
		}
	}

	if sm.Dirty {
		cmds = append(cmds, m.passwordInput.Focus())
	}

	return tea.Batch(cmds...)
}
