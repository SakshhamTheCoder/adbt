package adb

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Device struct {
	Serial  string
	Model   string
	State   string
	Android string
}

type DevicesLoadedMsg struct {
	Devices []Device
	Error   error
}

type PairWirelessResultMsg struct {
	Error error
}

func ListDevices() tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand("", "devices")
		if err != nil {
			return DevicesLoadedMsg{
				Error: fmt.Errorf("failed to list devices: %w", err),
			}
		}

		lines := ParseLines(out)

		if len(lines) > 0 && strings.Contains(
			strings.ToLower(lines[0]),
			"list of devices",
		) {
			lines = lines[1:]
		}

		devices := make([]Device, 0, len(lines))

		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			device := Device{
				Serial: parts[0],
				State:  parts[1],
			}

			if device.State == "device" {
				device.Model, _ = GetProperty(device.Serial, "ro.product.model")
				device.Android, _ = GetProperty(device.Serial, "ro.build.version.release")
			}

			devices = append(devices, device)
		}

		return DevicesLoadedMsg{Devices: devices}
	}
}

func PairWirelessCmd(addr, port, pin string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand("", "pair", addr+":"+port, pin)
		if err != nil {
			return PairWirelessResultMsg{
				Error: fmt.Errorf("failed to pair with device %s: %w", addr, err),
			}
		}

		output := string(out)

		if !strings.Contains(output, "Successfully paired") {
			return PairWirelessResultMsg{
				Error: fmt.Errorf(
					"failed to pair with device %s: %s",
					addr,
					strings.TrimSpace(output),
				),
			}
		}

		return PairWirelessResultMsg{Error: nil}
	}
}

func (d *Device) DisplayName() string {
	if d.Model != "" {
		return fmt.Sprintf("%s (%s)", d.Model, d.Serial)
	}
	return d.Serial
}

func (d *Device) IsConnected() bool {
	return d.State == "device"
}
