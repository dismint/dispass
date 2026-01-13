package uconst

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var (
	PasswordChar       = '▪'
	PaginatorDotString = "▪"
)

// https://lospec.com/palette-list/lost-century
const (
	lostCentury1  = "#d1b187"
	lostCentury2  = "#c77b58"
	lostCentury3  = "#ae5d40"
	lostCentury4  = "#79444a"
	lostCentury5  = "#4b3d44"
	lostCentury6  = "#ba9158"
	lostCentury7  = "#927441"
	lostCentury8  = "#4d4539"
	lostCentury9  = "#77743b"
	lostCentury10 = "#b3a555"
	lostCentury11 = "#d2c9a5"
	lostCentury12 = "#8caba1"
	lostCentury13 = "#4b726e"
	lostCentury14 = "#574852"
	lostCentury15 = "#847875"
	lostCentury16 = "#ab9b8e"
)

var (
	SymbolStyle              lipgloss.Style
	TextStyle                lipgloss.Style
	HelpKeyStyle             lipgloss.Style
	HelpDescStyle            lipgloss.Style
	HelpSeparatorStyle       lipgloss.Style
	MessageLevelErrorStyle   lipgloss.Style
	MessageLevelSuccessStyle lipgloss.Style
	MessageLevelNotifStyle   lipgloss.Style
	BorderColor              lipgloss.AdaptiveColor
	HelpStyles               help.Styles
)

var (
	ViewStyle         = lipgloss.NewStyle().Padding(1, 2).Width(50)
	ViewportViewStyle = lipgloss.NewStyle().Padding(0, 1).Width(44).
				Border(lipgloss.RoundedBorder()).BorderForeground(BorderColor)
)

func NewTextInput(prompt string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.PromptStyle = SymbolStyle
	ti.Cursor.Style = SymbolStyle
	ti.TextStyle = TextStyle
	ti.Width = 31
	return ti
}

func TruncAndPadListElem(text string) string {
	truncText := truncate.StringWithTail(text, 20, "…")
	return TextStyle.Width(21).Render(truncText)
}
