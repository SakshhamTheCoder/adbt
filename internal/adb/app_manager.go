package adb

import (
	"fmt"
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

type AppActionResultMsg struct {
	Action string
}

type AppActionErrorMsg struct {
	Action string
	Error  error
}

func ParseApps(output []byte) []App {
	lines := ParseLines(output)
	apps := make([]App, 0, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(line, "package:") {
			continue
		}

		line = strings.TrimPrefix(line, "package:")

		parts := strings.LastIndex(line, "=")
		if parts == -1 {
			continue
		}

		apkPath := line[:parts]
		pkgName := line[parts+1:]

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

func LaunchAppCmd(serial, pkg string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(serial, "shell", "cmd", "package", "resolve-activity", "--brief", "-a", "android.intent.action.MAIN", "-c", "android.intent.category.LAUNCHER", pkg)
		if err != nil {
			out, err = ExecuteCommand(serial, "shell", "cmd", "package", "resolve-activity", "--brief", pkg)
			if err != nil {
				return AppActionErrorMsg{Action: "launch", Error: fmt.Errorf("failed to find activity: %w", err)}
			}
		}

		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var component string
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if strings.Contains(line, "/") {
				component = line
				break
			}
		}

		if component == "" {
			return AppActionErrorMsg{Action: "launch", Error: fmt.Errorf("no launchable activity found for %s", pkg)}
		}

		_, err = ExecuteCommand(serial, "shell", "am", "start", "-n", component)
		if err != nil {
			return AppActionErrorMsg{Action: "launch", Error: err}
		}
		return AppActionResultMsg{Action: "launch"}
	}
}

func ForceStopAppCmd(serial, pkg string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"am",
			"force-stop",
			pkg,
		)
		if err != nil {
			return AppActionErrorMsg{Action: "force stop", Error: err}
		}
		return AppActionResultMsg{Action: "force stop"}
	}
}

func UninstallAppCmd(serial, pkg string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"uninstall",
			pkg,
		)
		if err != nil {
			return AppActionErrorMsg{Action: "uninstall", Error: err}
		}
		return AppActionResultMsg{Action: "uninstall"}
	}
}

func ClearAppDataCmd(serial, pkg string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"pm",
			"clear",
			pkg,
		)
		if err != nil {
			return AppActionErrorMsg{Action: "clear data", Error: err}
		}
		return AppActionResultMsg{Action: "clear data"}
	}
}

func isSystemApp(apkPath string) bool {
	return strings.HasPrefix(apkPath, "/system") ||
		strings.HasPrefix(apkPath, "/vendor") ||
		strings.HasPrefix(apkPath, "/product")
}
