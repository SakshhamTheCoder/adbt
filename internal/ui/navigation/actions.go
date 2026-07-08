package navigation

import (
	"github.com/SakshhamTheCoder/adbt/internal/state"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen string

const (
	ScreenDashboard   Screen = "dashboard"
	ScreenDevices     Screen = "devices"
	ScreenDeviceInfo  Screen = "device_info"
	ScreenLogcat      Screen = "logcat"
	ScreenApps        Screen = "apps"
	ScreenFiles       Screen = "files"
	ScreenPerfMonitor Screen = "perf_monitor"
	ScreenIntents     Screen = "intents"
	ScreenPorts       Screen = "ports"
	ScreenInput       Screen = "input"
)

type ScreenMeta struct {
	Title         string
	RequireDevice bool
}

// Registry is the single source of truth for the app's screens: their display
// titles and whether they need a connected device.
var Registry = map[Screen]ScreenMeta{
	ScreenDashboard:   {Title: "Dashboard"},
	ScreenDevices:     {Title: "Device Selection"},
	ScreenDeviceInfo:  {Title: "Device Info", RequireDevice: true},
	ScreenLogcat:      {Title: "Logcat", RequireDevice: true},
	ScreenApps:        {Title: "Apps", RequireDevice: true},
	ScreenFiles:       {Title: "Files", RequireDevice: true},
	ScreenPerfMonitor: {Title: "Performance", RequireDevice: true},
	ScreenIntents:     {Title: "Intents", RequireDevice: true},
	ScreenPorts:       {Title: "Ports", RequireDevice: true},
	ScreenInput:       {Title: "Input", RequireDevice: true},
}

// Navigate switches to the given screen, redirecting to device selection when
// the target needs a device and none is currently selected.
func Navigate(screen Screen, st *state.AppState) tea.Cmd {
	target := screen
	if Registry[screen].RequireDevice && !st.HasDevice() {
		target = ScreenDevices
	}
	return func() tea.Msg {
		return SwitchScreenMsg{Screen: target}
	}
}
