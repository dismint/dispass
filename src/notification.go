package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MessageLevel int

const (
	MessageLevelError MessageLevel = iota
	MessageLevelSuccess
	MessageLevelNotif
)

type showNotificationMsg string
type clearNotificationMsg struct{}

func showNotification(message string, messageLevel MessageLevel) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			var messageStyle lipgloss.Style
			switch messageLevel {
			case MessageLevelError:
				messageStyle = messageLevelErrorStyle
			case MessageLevelSuccess:
				messageStyle = messageLevelSuccessStyle
			case MessageLevelNotif:
				messageStyle = messageLevelNotifStyle
			}
			return showNotificationMsg(messageStyle.Render(message))
		},
		tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return clearNotificationMsg{}
		}),
	)
}
