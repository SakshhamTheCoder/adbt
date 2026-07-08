package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := adb.EnsureAvailable(); err != nil {
		fmt.Fprintln(os.Stderr, "adb not found in PATH.")
		fmt.Fprintln(os.Stderr, "Install Android platform-tools and try again:")
		fmt.Fprintln(os.Stderr, "  https://developer.android.com/tools/releases/platform-tools")
		os.Exit(1)
	}

	p := tea.NewProgram(
		ui.NewApp(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
