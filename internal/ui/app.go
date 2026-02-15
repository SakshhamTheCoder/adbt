package ui

import (
	"adbt/internal/state"
	"adbt/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	state         *state.AppState
	currentScreen tea.Model
	screenName    string
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
	return a.currentScreen.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.state.Width = msg.Width
		a.state.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit

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

	case screens.SwitchScreenMsg:
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

	default:
		return a, nil
	}

	a.currentScreen = newScreen
	a.screenName = name
	return a, newScreen.Init()
}

func (a *App) View() string {
	if a.state.Width == 0 {
		return "Initializing..."
	}
	return a.currentScreen.View()
}
