package changemaster

import (
	"fmt"

	"github.com/dismint/dispass/internal/uconst"
)

func (m *Model) View() string {
	view := fmt.Sprintf("%v\n\n%v\n",
		m.helpModel.View(m.keyMap),
		m.passwordInput.View(),
	)
	if m.confirming {
		view += fmt.Sprintf("%v\n",
			m.confirmPasswordInput.View(),
		)
	}
	return uconst.ViewStyle.Render(view)
}
