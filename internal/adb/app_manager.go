package adb

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	PackageName string
	APKPath     string
	IsSystem    bool
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
		if !strings.HasPrefix(line, "package:") {
			continue
		}

		line = strings.TrimPrefix(line, "package:")

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		apkPath := parts[0]
		pkgName := parts[1]

		apps = append(apps, App{
			PackageName: pkgName,
			APKPath:     apkPath,
			IsSystem:    isSystemApp(apkPath),
		})
	}

	sort.Slice(apps, func(i, j int) bool {
		return apps[i].PackageName < apps[j].PackageName
	})

	return apps
}

func ListAppsCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(
			serial,
			"shell",
			"pm",
			"list",
			"packages",
			"-f",
		)
		if err != nil {
			return AppsLoadErrorMsg{Error: err}
		}

		return AppsLoadedMsg{
			Apps: ParseApps(out),
		}
	}
}

func isSystemApp(apkPath string) bool {
	return strings.HasPrefix(apkPath, "/system") ||
		strings.HasPrefix(apkPath, "/vendor") ||
		strings.HasPrefix(apkPath, "/product")
}
