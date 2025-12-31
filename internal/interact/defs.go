package interact

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

type SearchKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}
type NavKeyMap struct {
	Quit   key.Binding
	Search key.Binding
	Clear  key.Binding
	Nav    key.Binding
	Copy   key.Binding
	Edit   key.Binding
	New    key.Binding
	Del    key.Binding
}
type ViewportKeyMap struct {
	Quit key.Binding
	Back key.Binding
	Save key.Binding
	Next key.Binding
	Prev key.Binding
}

func (k SearchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Confirm}
}
func (k NavKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Search, k.Clear, k.Nav, k.Copy, k.Edit, k.New, k.Del}
}
func (k ViewportKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Back, k.Save, k.Next, k.Prev}
}

func (k SearchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Confirm},
	}
}
func (k NavKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Search, k.Clear, k.Nav},
		{k.Copy, k.Edit, k.New, k.Del},
	}
}
func (k ViewportKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Back, k.Save},
		{k.Next, k.Prev},
	}
}

var searchKeyMap = SearchKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter", "esc"),
		key.WithHelp("↵ / esc", "confirm"),
	),
}
var navKeyMap = NavKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Search: key.NewBinding(
		key.WithKeys("s", "/"),
		key.WithHelp("s /", "search"),
	),
	Clear: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear"),
	),
	Nav: key.NewBinding(
		key.WithKeys("left", "up", "right", "down", "h", "k", "l", "j"),
		key.WithHelp("←↑→↓", "nav"),
	),
	Copy: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "copy"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Del: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}
var viewportKeyMap = ViewportKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "save"),
	),
	Next: key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("↓ tab", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "prev"),
	),
}

type Mode int

const (
	ModeSearch Mode = iota
	ModeNav
	ModeViewport
)

type Model struct {
	keyMap    help.KeyMap
	helpModel help.Model

	mode Mode

	viewportSourceInput   textinput.Model
	viewportUsernameInput textinput.Model
	viewportPasswordInput textinput.Model
	viewportUUID          string

	lastQuery       string
	keyInput        textinput.Model
	resultPaginator paginator.Model
	resultLocOnPage int
	topIDs          []string
}

func Initial() Model {
	keyInput := uconst.NewTextInput("search/")
	keyInput.ShowSuggestions = true

	viewportSourceInput := uconst.NewTextInput("Source    ")
	viewportUsernameInput := uconst.NewTextInput("Username  ")
	viewportPasswordInput := uconst.NewTextInput("Password  ")
	viewportPasswordInput.EchoMode = textinput.EchoPassword
	viewportPasswordInput.EchoCharacter = uconst.PasswordChar

	resultPaginator := paginator.New()
	resultPaginator.Type = paginator.Dots
	resultPaginator.PerPage = 10
	resultPaginator.ActiveDot = uconst.SymbolStyle.Render(uconst.PaginatorDotString)
	resultPaginator.InactiveDot = uconst.TextStyle.Render(uconst.PaginatorDotString)

	helpModel := help.New()
	helpModel.Styles = uconst.HelpStyles
	helpModel.ShowAll = true

	return Model{
		keyMap:    navKeyMap,
		helpModel: helpModel,

		mode: ModeNav,

		viewportSourceInput:   viewportSourceInput,
		viewportUsernameInput: viewportUsernameInput,
		viewportPasswordInput: viewportPasswordInput,
		// viewportUUID

		// lastQuery
		keyInput:        keyInput,
		resultPaginator: resultPaginator,
		// resultLocOnPage
		topIDs: make([]string, 0),
	}
}

func (m *Model) getSelectedCredInfo(sm *state.Model) (state.CredInfo, string, bool) {
	if len(m.topIDs) == 0 {
		return state.CredInfo{}, "", false
	}

	start, _ := m.resultPaginator.GetSliceBounds(len(m.topIDs))
	id := m.topIDs[start+m.resultLocOnPage]
	credInfo, exists := sm.KeyToCredInfo[id]
	return credInfo, id, exists
}
