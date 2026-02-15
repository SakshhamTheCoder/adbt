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

/* ---------- extended device details ---------- */

type DeviceDetails struct {
	BatteryLevel  string
	BatteryStatus string
	StorageUsed   string
	StorageTotal  string
	ScreenSize    string
	ScreenDensity string
	IPAddress     string
}

type DeviceDetailsMsg struct {
	Details DeviceDetails
	Error   error
}

func FetchDeviceDetailsCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		var d DeviceDetails

		out, err := ExecuteCommand(serial, "shell", "dumpsys", "battery")
		if err == nil {
			d.BatteryLevel, d.BatteryStatus = parseBattery(string(out))
		}

		out, err = ExecuteCommand(serial, "shell", "df", "/data")
		if err == nil {
			d.StorageUsed, d.StorageTotal = parseStorage(string(out))
		}

		out, err = ExecuteCommand(serial, "shell", "wm", "size")
		if err == nil {
			d.ScreenSize = parseWmOutput(string(out))
		}

		out, err = ExecuteCommand(serial, "shell", "wm", "density")
		if err == nil {
			d.ScreenDensity = parseWmOutput(string(out))
		}

		out, err = ExecuteCommand(serial, "shell", "ip", "route")
		if err == nil {
			d.IPAddress = parseIPAddress(string(out))
		}

		return DeviceDetailsMsg{Details: d}
	}
}

func parseBattery(output string) (level, status string) {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "level:") {
			level = strings.TrimSpace(strings.TrimPrefix(line, "level:")) + "%"
		}

		if strings.HasPrefix(line, "status:") {
			code := strings.TrimSpace(strings.TrimPrefix(line, "status:"))
			switch code {
			case "1":
				status = "Unknown"
			case "2":
				status = "Charging"
			case "3":
				status = "Discharging"
			case "4":
				status = "Not charging"
			case "5":
				status = "Full"
			default:
				status = code
			}
		}
	}
	return
}

func parseStorage(output string) (used, total string) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return
	}

	fields := strings.Fields(lines[1])
	if len(fields) >= 3 {
		total = formatStorageBlocks(fields[1])
		used = formatStorageBlocks(fields[2])
	}
	return
}

func formatStorageBlocks(blocksStr string) string {
	var n int64
	_, err := fmt.Sscanf(blocksStr, "%d", &n)
	if err != nil {
		return blocksStr
	}

	bytes := n * 1024

	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	default:
		return fmt.Sprintf("%d KB", n)
	}
}

func parseWmOutput(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, ":"); idx != -1 {
			return strings.TrimSpace(line[idx+1:])
		}
	}
	return strings.TrimSpace(output)
}

func parseIPAddress(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "src") {
			fields := strings.Fields(line)
			for i, f := range fields {
				if f == "src" && i+1 < len(fields) {
					return fields[i+1]
				}
			}
		}
	}
	return "N/A"
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
