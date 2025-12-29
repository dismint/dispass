package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// https://lospec.com/palette-list/seafoam
const seafoam1 = "#37364e"
const seafoam2 = "#355d69"
const seafoam3 = "#6aae9d"
const seafoam4 = "#b9d4b4"
const seafoam5 = "#f4e9d4"
const seafoam6 = "#d0baa9"
const seafoam7 = "#9e8e91"
const seafoam8 = "#5b4a68"

const lpred = "#f8858b"

var (
	symbolStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam1,
		Dark:  seafoam3,
	})
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam1,
		Dark:  seafoam5,
	})
	helpKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam6,
	})
	helpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam7,
	})
	helpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam8,
	})
)

var helpStyles = help.Styles{
	ShortKey:       helpKeyStyle,
	ShortDesc:      helpDescStyle,
	ShortSeparator: helpSeparatorStyle,
	FullKey:        helpKeyStyle,
	FullDesc:       helpDescStyle,
	FullSeparator:  helpSeparatorStyle,
}

var (
	finalWrapStyle        = lipgloss.NewStyle().Padding(1, 2)
	finalWrapStyleBounded = finalWrapStyle.Width(50)
	viewportStyle         = lipgloss.NewStyle().
				Padding(0, 1).
				Width(44).
				Border(lipgloss.RoundedBorder())
	messageBaseStyle       = lipgloss.NewStyle().Padding(0, 1)
	messageLevelErrorStyle = messageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  lpred,
	})
	messageLevelSuccessStyle = messageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam4,
	})
	messageLevelNotifStyle = messageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam5,
	})
)

func newTextInput(prompt string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.PromptStyle = symbolStyle
	ti.Cursor.Style = symbolStyle
	ti.TextStyle = textStyle
	ti.Width = 31
	return ti
}

func truncAndPadListElem(text string) string {
	truncText := truncate.StringWithTail(text, 20, "…")
	return lipgloss.NewStyle().Width(21).Render(truncText)
}

const logFileName = "dp.log"
const dataFileName = "dp.dat"
const bleveDirName = "index"
