package adb

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type DeviceActionResultMsg struct {
	Action string
}

type DeviceActionErrorMsg struct {
	Action string
	Error  error
}

func RebootCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "reboot")
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "reboot",
				Error:  err,
			}
		}
		return DeviceActionResultMsg{Action: "reboot"}
	}
}

func RebootRecoveryCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "reboot recovery")
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "reboot recovery",
				Error:  err,
			}
		}
		return DeviceActionResultMsg{Action: "reboot recovery"}
	}
}

func RebootBootloaderCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(serial, "reboot bootloader")
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "reboot bootloader",
				Error:  err,
			}
		}
		return DeviceActionResultMsg{Action: "reboot bootloader"}
	}
}

func StartScrcpyCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("scrcpy", "-s", serial)
		err := cmd.Start()

		if err != nil {
			return DeviceActionErrorMsg{
				Action: "scrcpy",
				Error:  err,
			}
		}
		return DeviceActionResultMsg{Action: "scrcpy"}
	}
}
