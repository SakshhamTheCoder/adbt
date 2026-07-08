package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// Populated by goreleaser at build time via -X ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.BoolVar(showVersion, "v", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("adbt %s (commit %s, built %s)\n", version, commit, date)
		return
	}

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
