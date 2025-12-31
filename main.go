package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/master"
	"github.com/dismint/dispass/internal/uconst"
)

func main() {
	logFd, err := os.OpenFile(
		uconst.LogFileName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	defer logFd.Close()
	log.SetOutput(logFd)

	if _, err := tea.NewProgram(master.Initial()).Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}
