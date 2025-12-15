package screens

import (
	"strings"

	"adbt/internal/adb"
	"adbt/internal/state"
	"adbt/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type Logcat struct {
	state   *state.AppState
	lines   []string
	session *adb.LogcatSession
	running bool
}

func NewLogcat(state *state.AppState) *Logcat {
	return &Logcat{state: state}
}

func (l *Logcat) Init() tea.Cmd {
	if !l.state.HasDevice() {
		return nil
	}
	l.running = true
	return adb.StartLogcatCmd(l.state.DeviceSerial())
}

func (l *Logcat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case adb.LogcatStartedMsg:
		l.session = msg.Session
		return l, adb.NextLogcatLineCmd(l.session)

	case adb.LogcatLineMsg:
		l.lines = append(l.lines, msg.Line)
		if len(l.lines) > 1000 {
			l.lines = l.lines[len(l.lines)-1000:]
		}
		if l.running {
			return l, adb.NextLogcatLineCmd(l.session)
		}

	case adb.LogcatStoppedMsg:
		l.running = false

	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			l.lines = nil
		case "s":
			l.running = !l.running
			if l.running {
				return l, adb.NextLogcatLineCmd(l.session)
			}
		}
	}

	return l, nil
}

func (l *Logcat) View() string {
	if !l.state.HasDevice() {
		return components.RenderLayout(l.state, components.LayoutProps{
			Title: "Logcat",
			Body:  components.StatusDisconnected.Render("No device selected"),
		})
	}

	var body strings.Builder
	for _, line := range l.lines {
		body.WriteString(line + "\n")
	}

	return components.RenderLayout(l.state, components.LayoutProps{
		Title: "Logcat",
		Body:  body.String(),
		Footer: components.Help("c", "clear") + "  " +
			components.Help("s", "start/stop"),
	})
}
