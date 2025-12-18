package screens

import (
	"strings"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type appAction struct {
	key   string
	label string
	cmd   func(string) tea.Cmd
}

type AppManager struct {
	state   *state.AppState
	actions []appAction
	apps    []adb.App
	loading bool
}

func NewAppManager(state *state.AppState) *AppManager {
	return &AppManager{
		state: state,
		actions: []appAction{
			{"l", "List Installed Apps", adb.ListAppsCmd},
		},
	}
}

func (a *AppManager) Init() tea.Cmd {
	return nil
}

func (a *AppManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !a.state.HasDevice() {
			return a, nil
		}
		for _, act := range a.actions {
			if msg.String() == act.key {
				a.loading = true
				return a, act.cmd(a.state.DeviceSerial())
			}
		}

	case adb.AppsLoadedMsg:
		a.loading = false
		a.apps = msg.Apps

	case adb.AppsLoadErrorMsg:
		a.loading = false

	}

	return a, nil
}

func (a *AppManager) View() string {
	if !a.state.HasDevice() {
		return components.RenderLayout(a.state, components.LayoutProps{
			Title: "Apps",
			Body:  components.StatusDisconnected.Render("No device selected"),
		})
	}

	var body strings.Builder

	if a.loading {
		body.WriteString("Loading apps...")
	} else {
		for _, app := range a.apps {
			body.WriteString(app.PackageName + "\n")
		}
	}

	return components.RenderLayout(a.state, components.LayoutProps{
		Title:  "Apps",
		Body:   body.String(),
		Footer: components.Help("l", "list apps") + "  " + components.Help("esc", "back"),
	})
}
