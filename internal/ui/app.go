package ui

import (
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"
	"github.com/SakshhamTheCoder/adbt/internal/ui/navigation"
	"github.com/SakshhamTheCoder/adbt/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	state         *state.AppState
	currentScreen tea.Model
	screenName    string
}

type LifecycleScreen interface {
	tea.Model
	Cleanup() tea.Cmd
}

func NewApp() *App {
	appState := state.New()

	return &App{
		state:         appState,
		currentScreen: screens.NewDashboard(appState),
		screenName:    "dashboard",
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(a.setAppTitle(), a.currentScreen.Init())
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.state.Width = msg.Width
		a.state.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Batch(a.cleanupCurrentScreen(), tea.Quit)

		case "esc":
			var cmd tea.Cmd
			a.currentScreen, cmd = a.currentScreen.Update(msg)
			if cmd != nil {
				return a, cmd
			}
			if a.screenName != "dashboard" {
				return a.switchScreen("dashboard")
			}
			return a, nil
		}

	case navigation.SwitchScreenMsg:
		return a.switchScreen(msg.Screen)
	}

	var cmd tea.Cmd
	a.currentScreen, cmd = a.currentScreen.Update(msg)
	return a, cmd
}

func (a *App) switchScreen(name string) (*App, tea.Cmd) {
	var newScreen tea.Model

	switch name {
	case "apps":
		newScreen = screens.NewAppManager(a.state)
	case "dashboard":
		newScreen = screens.NewDashboard(a.state)
	case "devices":
		newScreen = screens.NewDevices(a.state)
	case "device_info":
		newScreen = screens.NewDeviceInfo(a.state)
	case "files":
		newScreen = screens.NewFiles(a.state)
	case "logcat":
		newScreen = screens.NewLogcat(a.state)
	case "perf_monitor":
		newScreen = screens.NewPerfMonitor(a.state)
	case "intents":
		newScreen = screens.NewIntents(a.state)
	case "ports":
		newScreen = screens.NewPorts(a.state)

	default:
		return a, nil
	}

	cleanupCmd := a.cleanupCurrentScreen()
	a.currentScreen = newScreen
	a.screenName = name
	return a, tea.Batch(cleanupCmd, a.setAppTitle(), newScreen.Init())
}

func (a *App) View() string {
	if a.state.Width == 0 {
		return "Initializing..."
	}
	return a.currentScreen.View()
}

func (a *App) cleanupCurrentScreen() tea.Cmd {
	screen, ok := a.currentScreen.(LifecycleScreen)
	if !ok {
		return nil
	}
	return screen.Cleanup()
}

func (a *App) setAppTitle() tea.Cmd {
	// ADBT explicitly requests its title on startup and screen switches. Some
	// terminals with shell integration may still override titles temporarily.
	return tea.SetWindowTitle(components.ShellTitle(a.state, screenDisplayTitle(a.screenName)))
}

func screenDisplayTitle(name string) string {
	switch name {
	case "apps":
		return "Apps"
	case "dashboard":
		return "Dashboard"
	case "devices":
		return "Device Selection"
	case "device_info":
		return "Device Info"
	case "files":
		return "Files"
	case "logcat":
		return "Logcat"
	case "perf_monitor":
		return "Performance"
	case "intents":
		return "Intents"
	case "ports":
		return "Ports"
	default:
		return name
	}
}
