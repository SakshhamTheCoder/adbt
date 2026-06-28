package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/SakshhamTheCoder/adbt/internal/adb"
	"github.com/SakshhamTheCoder/adbt/internal/state"
	"github.com/SakshhamTheCoder/adbt/internal/ui/components"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FilterAll AppFilter = iota
	FilterUser
	FilterSystem
)

type AppFilter int

var filterNames = []string{"All", "User", "System"}

type AppManager struct {
	state    *state.AppState
	apps     []adb.App
	filtered []adb.App // cache of apps after search/filter; rebuilt via applyFilter
	loading  bool
	cursor   int

	search components.SearchState

	filterType AppFilter

	viewport viewport.Model

	confirm components.ConfirmPrompt
	toast   components.Toast
	pending string

	installForm components.FormModal

	details  map[string]adb.AppDetails // lazily fetched, keyed by package
	selected map[string]bool           // packages marked for batch actions
}

type appDetailTickMsg struct{ pkg string }

func NewAppManager(state *state.AppState) *AppManager {
	return &AppManager{
		state:    state,
		viewport: viewport.New(0, 0),
		details:  map[string]adb.AppDetails{},
		selected: map[string]bool{},
	}
}

// scheduleDetail debounces detail fetches: it waits briefly, then fetches only if
// the selection still points at the same (uncached) package — so scrolling fast
// doesn't flood adb with dumpsys calls.
func (a *AppManager) scheduleDetail() tea.Cmd {
	app := a.selectedApp()
	if app == nil {
		return nil
	}
	if _, ok := a.details[app.PackageName]; ok {
		return nil
	}
	pkg := app.PackageName
	return tea.Tick(250*time.Millisecond, func(time.Time) tea.Msg {
		return appDetailTickMsg{pkg: pkg}
	})
}

func (a *AppManager) Init() tea.Cmd {
	if !a.state.HasDevice() {
		return nil
	}

	a.loading = true
	return adb.ListAppsCmd(a.state.DeviceSerial())
}

// applyFilter rebuilds the cached filtered list. Call it whenever the app list,
// search query, or filter type changes.
func (a *AppManager) applyFilter() {
	filtered := a.filtered[:0]
	lowerSearch := strings.ToLower(a.search.Query)

	for _, app := range a.apps {
		switch a.filterType {
		case FilterUser:
			if app.IsSystem {
				continue
			}
		case FilterSystem:
			if !app.IsSystem {
				continue
			}
		}

		if a.search.Query != "" {
			if !strings.Contains(strings.ToLower(app.PackageName), lowerSearch) {
				continue
			}
		}

		filtered = append(filtered, app)
	}

	a.filtered = filtered
}

func (a *AppManager) selectedApp() *adb.App {
	if len(a.filtered) == 0 || a.cursor >= len(a.filtered) {
		return nil
	}
	return &a.filtered[a.cursor]
}

// CapturingText keeps "q" out of the global quit handler while a form or the
// search box is active.
func (a *AppManager) CapturingText() bool {
	return a.search.Active || a.installForm.Visible
}

// detailBlock renders the lazily-fetched details for the highlighted app.
func (a *AppManager) detailBlock() string {
	app := a.selectedApp()
	if app == nil {
		return ""
	}

	tag := "User app"
	if app.IsSystem {
		tag = "System app"
	}

	info := "loading…"
	if d, ok := a.details[app.PackageName]; ok {
		var parts []string
		if d.VersionName != "" {
			v := d.VersionName
			if d.VersionCode != "" {
				v += " (" + d.VersionCode + ")"
			}
			parts = append(parts, v)
		}
		if d.Size != "" {
			parts = append(parts, d.Size)
		}
		if d.TargetSdk != "" {
			parts = append(parts, "SDK "+d.TargetSdk)
		}
		if len(parts) == 0 {
			info = "—"
		} else {
			info = strings.Join(parts, "  •  ")
		}
	}

	apkStyle := lipgloss.NewStyle().MaxWidth(max(a.state.Width-6, 20)).Foreground(components.FgMuted)

	var b strings.Builder
	b.WriteString("\n" + components.SectionTitle("Details") + "\n")
	b.WriteString(components.StatusMuted.Render("  "+tag+"  •  "+info) + "\n")
	b.WriteString(apkStyle.Render("  "+app.APKPath) + "\n")
	return b.String()
}

func (a *AppManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	a.toast.Update(msg)

	if a.installForm.Visible {
		switch msg := msg.(type) {
		case components.FormSubmitMsg:
			values := msg.Values
			a.installForm.Hide()
			if len(values) > 0 && values[0] != "" {
				var toastCmd tea.Cmd
				a.toast, toastCmd = components.ShowToast(
					"Installing APK...",
					false,
					2*time.Second,
				)
				return a, tea.Batch(
					toastCmd,
					adb.InstallAppCmd(a.state.DeviceSerial(), values[0]),
				)
			}
			return a, nil
		case components.FormCancelMsg:
			a.installForm.Hide()
			return a, nil
		}
		return a, a.installForm.Update(msg)
	}

	if a.confirm.Visible {
		switch msg.(type) {

		case components.ConfirmYesMsg:
			a.confirm.Hide()
			serial := a.state.DeviceSerial()

			if a.pending == "uninstall_selected" {
				var cmds []tea.Cmd
				for pkg, sel := range a.selected {
					if sel {
						cmds = append(cmds, adb.UninstallAppCmd(serial, pkg))
					}
				}
				a.selected = map[string]bool{}
				a.pending = ""
				return a, tea.Batch(cmds...)
			}

			app := a.selectedApp()
			if app == nil {
				return a, nil
			}
			switch a.pending {
			case "uninstall":
				return a, adb.UninstallAppCmd(serial, app.PackageName)
			case "clear data":
				return a, adb.ClearAppDataCmd(serial, app.PackageName)
			case "force_stop":
				return a, adb.ForceStopAppCmd(serial, app.PackageName)
			}
			return a, nil

		case components.ConfirmNoMsg:
			a.confirm.Hide()
			a.pending = ""
			return a, tea.Batch()
		}

		return a, a.confirm.Update(msg)
	}

	switch msg := msg.(type) {

	case adb.AppsLoadedMsg:
		a.loading = false
		a.apps = msg.Apps
		a.applyFilter()
		a.cursor = 0
		a.gotoTop()
		return a, a.scheduleDetail()

	case appDetailTickMsg:
		app := a.selectedApp()
		if app == nil || app.PackageName != msg.pkg {
			return a, nil
		}
		if _, ok := a.details[app.PackageName]; ok {
			return a, nil
		}
		return a, adb.AppDetailsCmd(a.state.DeviceSerial(), app.PackageName, app.APKPath)

	case adb.AppDetailsMsg:
		if msg.Error == nil {
			a.details[msg.Pkg] = msg.Details
		}
		return a, nil

	case adb.AppsLoadErrorMsg:
		a.loading = false
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			"Failed to load apps",
			true,
			3*time.Second,
		)
		return a, cmd

	case adb.AppActionResultMsg:
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			msg.Action+" successful",
			false,
			2*time.Second,
		)
		if msg.Action == "uninstall" || msg.Action == "install" {
			return a, tea.Batch(
				cmd,
				adb.ListAppsCmd(a.state.DeviceSerial()),
			)
		}
		return a, cmd

	case adb.AppActionErrorMsg:
		var cmd tea.Cmd
		a.toast, cmd = components.ShowToast(
			msg.Action+" failed: "+msg.Error.Error(),
			true,
			3*time.Second,
		)
		return a, cmd

	case tea.KeyMsg:
		if a.search.Active {
			before := a.search.Query
			a.search.HandleKey(msg)
			if a.search.Query != before {
				a.applyFilter()
				a.cursor = 0
				a.gotoTop()
				return a, tea.Batch(consumeKeyCmd(), a.scheduleDetail())
			}
			return a, consumeKeyCmd()
		}

		filtered := a.filtered

		switch msg.String() {
		case "up", "k":
			if a.cursor > 0 {
				a.cursor--
				a.ensureCursorVisible()
			}
			return a, a.scheduleDetail()

		case "down", "j":
			if a.cursor < len(filtered)-1 {
				a.cursor++
				a.ensureCursorVisible()
			}
			return a, a.scheduleDetail()

		case " ", "space":
			if app := a.selectedApp(); app != nil {
				if a.selected[app.PackageName] {
					delete(a.selected, app.PackageName)
				} else {
					a.selected[app.PackageName] = true
				}
			}

		case "e":
			if app := a.selectedApp(); app != nil {
				var toastCmd tea.Cmd
				a.toast, toastCmd = components.ShowToast("Extracting APK...", false, 2*time.Second)
				return a, tea.Batch(
					toastCmd,
					adb.ExtractApkCmd(a.state.DeviceSerial(), app.APKPath, app.PackageName),
				)
			}

		case "enter", "l":
			if app := a.selectedApp(); app != nil {
				return a, adb.LaunchAppCmd(
					a.state.DeviceSerial(),
					app.PackageName,
				)
			}

		case "s":
			if app := a.selectedApp(); app != nil {
				a.pending = "force_stop"
				a.confirm.Show("Force stop:\n" + app.PackageName)
			}

		case "u":
			if len(a.selected) > 0 {
				a.pending = "uninstall_selected"
				a.confirm.Show(fmt.Sprintf("Uninstall %d selected app(s)?", len(a.selected)))
				return a, nil
			}
			if app := a.selectedApp(); app != nil {
				if app.IsSystem {
					var cmd tea.Cmd
					a.toast, cmd = components.ShowToast(
						"Cannot uninstall system app",
						true,
						2*time.Second,
					)
					return a, cmd
				}
				a.pending = "uninstall"
				a.confirm.Show("Uninstall:\n" + app.PackageName)
			}

		case "x":
			if app := a.selectedApp(); app != nil {
				a.pending = "clear data"
				a.confirm.Show("Clear data:\n" + app.PackageName)
			}

		case "/":
			a.search.Start()

		case "r":
			if a.state.HasDevice() {
				a.loading = true
				a.cursor = 0
				a.gotoTop()
				return a, adb.ListAppsCmd(a.state.DeviceSerial())
			}

		case "i":
			a.installForm.Show("Install APK", []components.FormField{
				{Label: "APK Path", Value: ""},
			})

		case "right":
			a.filterType = (a.filterType + 1) % 3
			a.applyFilter()
			a.cursor = 0
			a.gotoTop()
			return a, a.scheduleDetail()

		case "left":
			a.filterType = (a.filterType + 2) % 3
			a.applyFilter()
			a.cursor = 0
			a.gotoTop()
			return a, a.scheduleDetail()

		case "esc":
			if len(a.selected) > 0 {
				a.selected = map[string]bool{}
				return a, consumeKeyCmd()
			}
			if a.search.Query != "" {
				a.search.Clear()
				a.applyFilter()
				a.cursor = 0
				a.gotoTop()
				return a, consumeKeyCmd()
			}

		default:
			return a, a.updateViewport(msg)
		}
	}

	return a, nil
}

func (a *AppManager) View() string {
	if !a.state.HasDevice() {
		return components.RenderNoDevice(a.state, "Apps")
	}

	var staticContent strings.Builder

	staticContent.WriteString("  ")
	for i, name := range filterNames {
		if AppFilter(i) == a.filterType {
			staticContent.WriteString(components.TabActiveStyle.Render(name))
		} else {
			staticContent.WriteString(components.TabInactiveStyle.Render(name))
		}
		if i < len(filterNames)-1 {
			staticContent.WriteString(" ")
		}
	}
	if len(a.selected) > 0 {
		staticContent.WriteString(
			"  " + components.HelpKeyStyle.Render(fmt.Sprintf("%d selected", len(a.selected))),
		)
	}
	staticContent.WriteString("\n")

	if a.search.Active {
		staticContent.WriteString(
			components.HelpKeyStyle.Render("search: ") + a.search.Query + "▌\n",
		)
	} else if a.search.Query != "" {
		staticContent.WriteString(
			components.StatusMuted.Render("filter: \""+a.search.Query+"\"") + "\n",
		)
	}

	staticContent.WriteString(a.detailBlock())

	maxWidth := max(a.state.Width-8, 20)
	truncStyle := lipgloss.NewStyle().MaxWidth(maxWidth)

	var scrollableContent strings.Builder

	if a.loading {
		scrollableContent.WriteString(components.StatusMuted.Render("Loading apps..."))
	} else {
		filtered := a.filtered

		if len(filtered) == 0 {
			scrollableContent.WriteString(components.StatusMuted.Render("No apps found"))
		} else {
			for i, app := range filtered {
				prefix := "  "
				if i == a.cursor {
					prefix = "› "
				}

				tag := components.StatusMuted.Render("[U]")
				if app.IsSystem {
					tag = components.StatusMuted.Render("[S]")
				}

				mark := " "
				if a.selected[app.PackageName] {
					mark = components.StatusConnected.Render("✓")
				}

				var line string
				if i == a.cursor {
					line = fmt.Sprintf(
						"%s%s %s %s",
						prefix,
						mark,
						tag,
						components.ListItemSelectedStyle.Render(app.PackageName),
					)
				} else {
					line = fmt.Sprintf(
						"%s%s %s %s",
						prefix,
						mark,
						tag,
						components.ListItemStyle.Render(app.PackageName),
					)
				}

				scrollableContent.WriteString(truncStyle.Render(line) + "\n")
			}
		}
	}

	var footer string
	if a.search.Active {
		footer = components.JoinHelp([2]string{"enter", "apply"}, [2]string{"esc", "cancel"})
	} else if a.search.Query != "" {
		footer = components.JoinHelp(
			[2]string{"↑/↓", "navigate"},
			[2]string{"enter", "launch"},
			[2]string{"/", "search"},
			[2]string{"esc", "clear filter"},
		)
	} else {
		footer = components.JoinHelp(
			[2]string{"↑/↓", "navigate"},
			[2]string{"space", "select"},
			[2]string{"enter", "launch"},
			[2]string{"i", "install"},
			[2]string{"e", "extract"},
			[2]string{"s", "stop"},
			[2]string{"u", "uninstall"},
			[2]string{"x", "clear"},
			[2]string{"←/→", "filter"},
			[2]string{"/", "search"},
			[2]string{"r", "reload"},
			[2]string{"esc", "back"},
		)
	}

	rendered := components.RenderLayoutWithScrollableSection(a.state, components.LayoutWithScrollProps{
		Title:             "Apps",
		StaticContent:     staticContent.String(),
		ScrollableContent: scrollableContent.String(),
		Footer:            footer,
		Viewport:          &a.viewport,
	})

	if a.installForm.Visible {
		rendered = components.RenderFormOverlay(rendered, a.installForm, a.state)
	}

	if a.confirm.Visible {
		rendered = components.RenderOverlay(rendered, a.confirm.View(), a.state)
	}

	if a.toast.Visible {
		rendered = components.RenderOverlay(rendered, a.toast.View(), a.state)
	}

	return rendered
}

func (a *AppManager) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	a.viewport, cmd = a.viewport.Update(msg)
	return cmd
}

func (a *AppManager) gotoTop() {
	a.viewport.GotoTop()
}

func (a *AppManager) ensureCursorVisible() {
	ensureViewportLineVisible(&a.viewport, a.cursor)
}

func consumeKeyCmd() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
