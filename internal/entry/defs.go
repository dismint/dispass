package entry

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dismint/dispass/internal/uconst"
)

type KeyMap struct {
	Quit  key.Binding
	Enter key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Enter}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Enter},
	}
}

var keyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "enter"),
	),
}

type Model struct {
	keyMap    KeyMap
	helpModel help.Model

	confirming bool

	passwordInput        textinput.Model
	confirmPasswordInput textinput.Model
}

func Initial() Model {
	passwordInput := textinput.New()
	passwordInput.Focus()
	passwordInput.CharLimit = -1
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = uconst.PasswordChar
	passwordInput.Prompt = uconst.PromptString
	passwordInput.Cursor.Style = uconst.SymbolStyle
	passwordInput.PromptStyle = uconst.SymbolStyle
	passwordInput.TextStyle = uconst.TextStyle

	confirmPasswordInput := textinput.New()
	confirmPasswordInput.CharLimit = -1
	confirmPasswordInput.EchoMode = textinput.EchoPassword
	confirmPasswordInput.EchoCharacter = uconst.PasswordChar
	confirmPasswordInput.Prompt = uconst.PromptString
	confirmPasswordInput.Cursor.Style = uconst.SymbolStyle
	confirmPasswordInput.PromptStyle = uconst.SymbolStyle
	confirmPasswordInput.TextStyle = uconst.TextStyle

	helpModel := help.New()
	helpModel.Styles = uconst.HelpStyles

	return Model{
		keyMap:    keyMap,
		helpModel: helpModel,

		confirming: false,

		passwordInput:        passwordInput,
		confirmPasswordInput: confirmPasswordInput,
	}
}
