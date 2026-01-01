package state

import (
	"time"

	"github.com/blevesearch/bleve"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dismint/dispass/internal/uconst"
)

type Screen int

const (
	EntryScreen Screen = iota
	InteractScreen
	ChangeMasterScreen
)

type MessageLevel int

const (
	MessageLevelError MessageLevel = iota
	MessageLevelSuccess
	MessageLevelNotif
)

type CredInfo struct {
	Source   string
	Username string
	Password string
}

type ShowNotificationMsg string
type ClearNotificationMsg struct{}

func NotificationMsg(message string, messageLevel MessageLevel) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			var messageStyle lipgloss.Style
			switch messageLevel {
			case MessageLevelError:
				messageStyle = uconst.MessageLevelErrorStyle
			case MessageLevelSuccess:
				messageStyle = uconst.MessageLevelSuccessStyle
			case MessageLevelNotif:
				messageStyle = uconst.MessageLevelNotifStyle
			}
			return ShowNotificationMsg(messageStyle.Render(message))
		},
		tea.Tick(1*time.Second, func(time.Time) tea.Msg {
			return ClearNotificationMsg{}
		}),
	)
}

type Model struct {
	Screen        Screen
	KeyToCredInfo map[string]CredInfo
	Secret        []byte
	Index         bleve.Index
	Notification  string
	Quitting      bool

	Dirty bool
}

func Initial() Model {
	return Model{
		Screen:        EntryScreen,
		KeyToCredInfo: make(map[string]CredInfo),
		// Secret
		// Index
		// Notification
		// Quitting

		Dirty: true,
	}
}

func (m *Model) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case ShowNotificationMsg:
		m.Notification = string(msg)
	case ClearNotificationMsg:
		m.Notification = ""
	}
}
