package main

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Screen int

const (
	EntryScreen Screen = iota
	SelectScreen
	ListScreen
)

type credInfo struct {
	Source   string
	Username string
	Password string
}

type model struct {
	screen Screen

	entryModel  EntryModel
	selectModel SelectModel

	keyToCredInfo map[string]credInfo
	secret        []byte

	notification string

	quitting bool
}

func initialModel() model {
	m := model{screen: EntryScreen}
	initialEntryModel(&m)
	initialSelectModel(&m)
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.entryModel.help.Width = msg.Width
	case showNotificationMsg:
		m.notification = string(msg)
	case clearNotificationMsg:
		m.notification = ""
	}

	switch m.screen {
	case EntryScreen:
		cmds = append(cmds, updateEntryModel(&m, msg))
	case SelectScreen:
		cmds = append(cmds, updateSelectModel(&m, msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var view string

	switch m.screen {
	case EntryScreen:
		view = viewEntryModel(&m)
	default:
		view = viewSelectModel(&m)
	}

	view += "\n" + m.notification

	return view
}

func main() {
	logFd, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	defer logFd.Close()
	log.SetOutput(logFd)

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}
