package interact

import (
	"fmt"

	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

func (m *Model) viewViewport(sm *state.Model) string {
	_, _, exists := m.getSelectedCredInfo(sm)
	if !exists && m.mode != ModeViewport {
		return "Feeling empty, create new credentials?"
	}
	return fmt.Sprintf("%v\n%v\n%v",
		m.viewportSourceInput.View(),
		m.viewportUsernameInput.View(),
		m.viewportPasswordInput.View(),
	)
}

func (m *Model) View(sm *state.Model) string {
	start, end := m.resultPaginator.GetSliceBounds(len(m.topIDs))

	var resultList string
	for locOnPage, topID := range m.topIDs[start:end] {
		prefix := " "
		if locOnPage == m.resultLocOnPage {
			prefix = uconst.SymbolStyle.Render(">")
		}
		resultList += fmt.Sprintf("%v %v %v\n",
			prefix,
			uconst.TruncAndPadListElem(sm.KeyToCredInfo[topID].Source),
			uconst.TruncAndPadListElem(sm.KeyToCredInfo[topID].Username),
		)
	}
	if resultList == "" {
		resultList = "No Results Found"
	}

	view := fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v\n\n%v",
		m.helpModel.View(m.keyMap),
		m.keyInput.View(),
		uconst.ViewportViewStyle.Render(m.viewViewport(sm)),
		m.resultPaginator.View(),
		resultList,
	)

	return uconst.ViewStyle.Render(view)
}
