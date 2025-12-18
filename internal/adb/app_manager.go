package adb

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	PackageName string
	Path        string
}

type AppsLoadedMsg struct {
	Apps []App
}

type AppsLoadErrorMsg struct {
	Error error
}

func ParseApps(output []byte) []App {
	lines := ParseLines(output)
	apps := make([]App, 0, len(lines))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var app App
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		app.PackageName = parts[1]
		pathParts := strings.SplitN(parts[0], ":", 2)
		if len(pathParts) != 2 {
			continue
		}
		app.Path = pathParts[1]
		apps = append(apps, app)
	}

	return apps
}

func ListAppsCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(serial, "shell", "pm", "list", "packages", "-f")
		if err != nil {
			return AppsLoadErrorMsg{Error: err}
		}

		return AppsLoadedMsg{
			Apps: ParseApps(out),
		}
	}
}
