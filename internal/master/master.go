package master

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dismint/dispass/internal/changemaster"
	"github.com/dismint/dispass/internal/entry"
	"github.com/dismint/dispass/internal/interact"
	"github.com/dismint/dispass/internal/state"
)

type Model struct {
	stateModel        state.Model
	entryModel        entry.Model
	interactModel     interact.Model
	changemasterModel changemaster.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func Initial() Model {
	return Model{
		stateModel:        state.Initial(),
		entryModel:        entry.Initial(),
		interactModel:     interact.Initial(),
		changemasterModel: changemaster.Initial(),
	}
}

func (m Model) screenUpdate(msg tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch m.stateModel.Screen {
	case state.EntryScreen:
		cmds = append(cmds, m.entryModel.Update(msg, &m.stateModel))
	case state.InteractScreen:
		cmds = append(cmds, m.interactModel.Update(msg, &m.stateModel))
	case state.ChangeMasterScreen:
		cmds = append(cmds, m.changemasterModel.Update(msg, &m.stateModel))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	m.stateModel.Update(msg)
	var cmd tea.Cmd
	m, cmd = m.screenUpdate(msg)
	cmds = append(cmds, cmd)

	if m.stateModel.Dirty {
		m, cmd = m.screenUpdate(msg)
		cmds = append(cmds, cmd)
		m.stateModel.Dirty = false
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.stateModel.Quitting {
		return ""
	}

	var view string

	switch m.stateModel.Screen {
	case state.EntryScreen:
		view = m.entryModel.View()
	case state.InteractScreen:
		view = m.interactModel.View(&m.stateModel)
	case state.ChangeMasterScreen:
		view = m.changemasterModel.View()
	}

	view += "\n" + m.stateModel.Notification

	return view
}
