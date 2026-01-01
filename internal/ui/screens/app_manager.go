package screens

import (
	"fmt"
	"strings"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type AppManager struct {
	state   *state.AppState
	apps    []adb.App
	loading bool
}

func NewAppManager(state *state.AppState) *AppManager {
	return &AppManager{
		state: state,
	}
}

func (a *AppManager) Init() tea.Cmd {
	if !a.state.HasDevice() {
		return nil
	}

	a.loading = true
	return adb.ListAppsCmd(a.state.DeviceSerial())
}

func (a *AppManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case adb.AppsLoadedMsg:
		a.loading = false
		a.apps = msg.Apps
		return a, nil

	case adb.AppsLoadErrorMsg:
		a.loading = false
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if a.state.HasDevice() {
				a.loading = true
				return a, adb.ListAppsCmd(a.state.DeviceSerial())
			}
		default:

			return a, components.UpdateViewport("Apps", msg)
		}
	}

	return a, nil
}

func (a *AppManager) View() string {
	if !a.state.HasDevice() {
		return components.RenderNoDevice(a.state, "Apps")
	}

	var body strings.Builder
	var scrollableContent strings.Builder

	body.WriteString(components.TitleStyle.Render("Installed Applications") + "\n\n")

	if a.loading {
		scrollableContent.WriteString(components.StatusMuted.Render("Loading apps..."))
	} else if len(a.apps) == 0 {
		scrollableContent.WriteString(components.StatusMuted.Render("No apps found"))
	} else {
		for _, app := range a.apps {
			tag := components.StatusMuted.Render("[U]")
			if app.IsSystem {
				tag = components.StatusMuted.Render("[S]")
			}

			line := fmt.Sprintf(
				"%s %s",
				tag,
				components.ListItemStyle.Render(app.PackageName),
			)

			scrollableContent.WriteString(line)
			scrollableContent.WriteString("\n")
		}
	}

	return components.RenderLayoutWithScrollableSection(a.state, components.LayoutWithScrollProps{
		Title:             "Apps",
		StaticContent:     body.String(),
		ScrollableContent: scrollableContent.String(),
		Footer: components.Help("↑/↓", "scroll") + "  " +
			components.Help("r", "reload") + "  " +
			components.Help("esc", "back"),
	})
}
