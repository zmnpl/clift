package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	wodb "github.com/zmnpl/clift/db"
	ui "github.com/zmnpl/clift/ui"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot find user home: %v", err)
	}

	// Init DB
	wodb.Init(filepath.Join(home, "Documents", "training.db"))

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}
