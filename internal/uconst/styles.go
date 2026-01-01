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
	SymbolStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam1,
		Dark:  seafoam3,
	})
	TextStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam1,
		Dark:  seafoam5,
	})
	HelpKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam6,
	})
	HelpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam7,
	})
	HelpSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam8,
	})
)

var HelpStyles = help.Styles{
	ShortKey:       HelpKeyStyle,
	ShortDesc:      HelpDescStyle,
	ShortSeparator: HelpSeparatorStyle,
	FullKey:        HelpKeyStyle,
	FullDesc:       HelpDescStyle,
	FullSeparator:  HelpSeparatorStyle,
}

var (
	ViewStyle         = lipgloss.NewStyle().Padding(1, 2).Width(50)
	ViewportViewStyle = lipgloss.NewStyle().Padding(0, 1).Width(44).
				Border(lipgloss.RoundedBorder())
	MessageBaseStyle       = lipgloss.NewStyle().Padding(0, 1)
	MessageLevelErrorStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  lpred,
	})
	MessageLevelSuccessStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam4,
	})
	MessageLevelNotifStyle = MessageBaseStyle.Foreground(lipgloss.AdaptiveColor{
		Light: seafoam7,
		Dark:  seafoam5,
	})
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
	return lipgloss.NewStyle().Width(21).Render(truncText)
}
