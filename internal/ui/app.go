package ui

import (
	"fmt"

	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"
	"github.com/SakshhamTheCoder/adbt/internal/ui/navigation"
	"github.com/SakshhamTheCoder/adbt/internal/ui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	state         *state.AppState
	currentScreen tea.Model
	screenName    navigation.Screen
}

type LifecycleScreen interface {
	tea.Model
	Cleanup() tea.Cmd
}

// TextCapturingScreen lets a screen signal that it is currently capturing text
// input, so global single-key shortcuts like "q" are forwarded to it instead of
// quitting the app.
type TextCapturingScreen interface {
	CapturingText() bool
}

func NewApp() *App {
	appState := state.New()

	return &App{
		state:         appState,
		currentScreen: screens.NewDashboard(appState),
		screenName:    navigation.ScreenDashboard,
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
		capturing := false
		if ts, ok := a.currentScreen.(TextCapturingScreen); ok {
			capturing = ts.CapturingText()
		}

		switch msg.String() {
		case "ctrl+c":
			return a, tea.Batch(a.cleanupCurrentScreen(), tea.Quit)

		case "q":
			if !capturing {
				return a, tea.Batch(a.cleanupCurrentScreen(), tea.Quit)
			}

		case "esc":
			var cmd tea.Cmd
			a.currentScreen, cmd = a.currentScreen.Update(msg)
			if cmd != nil {
				return a, cmd
			}
			if a.screenName != navigation.ScreenDashboard {
				return a.switchScreen(navigation.ScreenDashboard)
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

func (a *App) switchScreen(name navigation.Screen) (*App, tea.Cmd) {
	var newScreen tea.Model

	switch name {
	case navigation.ScreenApps:
		newScreen = screens.NewAppManager(a.state)
	case navigation.ScreenDashboard:
		newScreen = screens.NewDashboard(a.state)
	case navigation.ScreenDevices:
		newScreen = screens.NewDevices(a.state)
	case navigation.ScreenDeviceInfo:
		newScreen = screens.NewDeviceInfo(a.state)
	case navigation.ScreenFiles:
		newScreen = screens.NewFiles(a.state)
	case navigation.ScreenLogcat:
		newScreen = screens.NewLogcat(a.state)
	case navigation.ScreenPerfMonitor:
		newScreen = screens.NewPerfMonitor(a.state)
	case navigation.ScreenIntents:
		newScreen = screens.NewIntents(a.state)
	case navigation.ScreenPorts:
		newScreen = screens.NewPorts(a.state)
	case navigation.ScreenInput:
		newScreen = screens.NewInput(a.state)

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

	// Dimension guard
	if a.state.Width < components.MinScreenWidth || a.state.Height < components.MinScreenHeight {
		msg := fmt.Sprintf("Please increase dimension to at least %dx%d\n(Current: %dx%d)", components.MinScreenWidth, components.MinScreenHeight, a.state.Width, a.state.Height)
		return components.DimensionGuardStyle.
			Width(a.state.Width).
			Height(a.state.Height).
			Render(msg)
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
	title := navigation.Registry[a.screenName].Title
	if title == "" {
		title = string(a.screenName)
	}
	return tea.SetWindowTitle(components.ShellTitle(a.state, title))
}
