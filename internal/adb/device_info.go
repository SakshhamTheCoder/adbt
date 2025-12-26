package adb

import (
	"fmt"
	"os/exec"
	"strings"

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

func ToggleScreenCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		_, err := ExecuteCommand(
			serial,
			"shell",
			"input",
			"keyevent",
			"26",
		)
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "toggle screen",
				Error:  err,
			}
		}
		return DeviceActionResultMsg{Action: "toggle screen"}
	}
}

func ToggleWifiCmd(serial string) tea.Cmd {
	return func() tea.Msg {

		out, err := ExecuteCommand(serial, "shell", "dumpsys", "wifi")
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "check wifi state",
				Error:  err,
			}
		}

		output := string(out)
		enable := false

		if strings.Contains(output, "Wi-Fi is disabled") {
			enable = true
		} else if strings.Contains(output, "Wi-Fi is enabled") {
			enable = false
		} else {
			return DeviceActionErrorMsg{
				Action: "toggle wifi",
				Error:  fmt.Errorf("unable to determine wifi state"),
			}
		}

		action := "disable"
		if enable {
			action = "enable"
		}

		_, err = ExecuteCommand(
			serial,
			"shell",
			"svc",
			"wifi",
			action,
		)
		if err != nil {
			return DeviceActionErrorMsg{
				Action: "toggle wifi",
				Error:  err,
			}
		}

		return DeviceActionResultMsg{
			Action: "wifi " + action,
		}
	}
}
