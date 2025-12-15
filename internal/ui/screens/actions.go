package screens

import (
	"adbt/internal/state"

	tea "github.com/charmbracelet/bubbletea"
)

type Action string

const (
	ActionDevices    Action = "devices"
	ActionLogcat     Action = "logcat"
	ActionShell      Action = "shell"
	ActionApps       Action = "apps"
	ActionFiles      Action = "files"
	ActionDeviceInfo Action = "device_info"
)

func ResolveAction(action Action, state *state.AppState) tea.Cmd {
	switch action {

	case ActionDevices:
		return func() tea.Msg {
			return SwitchScreenMsg{Screen: "devices"}
		}

	case ActionLogcat:
		if !state.HasDevice() {
			return func() tea.Msg {
				return SwitchScreenMsg{Screen: "devices"}
			}
		}
		return func() tea.Msg {
			return SwitchScreenMsg{Screen: "logcat"}
		}

	case ActionShell, ActionApps, ActionFiles:
		if !state.HasDevice() {
			return func() tea.Msg {
				return SwitchScreenMsg{Screen: "devices"}
			}
		}
		return func() tea.Msg {
			return SwitchScreenMsg{Screen: string(action)}
		}

	case ActionDeviceInfo:
		if !state.HasDevice() {
			return func() tea.Msg {
				return SwitchScreenMsg{Screen: "devices"}
			}
		}
		return func() tea.Msg {
			return SwitchScreenMsg{Screen: "device_info"}
		}
	}

	return nil
}
