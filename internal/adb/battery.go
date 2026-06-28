package adb

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ThermalZone struct {
	Type string
	Temp string
}

type BatteryThermal struct {
	Level       string
	Status      string
	Temperature string
	Voltage     string
	Health      string
	Technology  string
	Plugged     string
	Zones       []ThermalZone
}

type BatteryThermalMsg struct {
	Data  BatteryThermal
	Error error
}

const maxThermalZones = 12

func FetchBatteryThermalCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		out, err := ExecuteCommand(serial, "shell", "dumpsys", "battery")
		if err != nil {
			return BatteryThermalMsg{Error: err}
		}
		data := parseBatteryFull(string(out))

		// Best effort — many devices restrict thermal_zone access.
		zonesOut, zErr := ExecuteCommand(serial, "shell",
			"for z in /sys/class/thermal/thermal_zone*; do echo $(cat $z/type):$(cat $z/temp); done")
		if zErr == nil {
			data.Zones = parseThermalZones(string(zonesOut))
		}

		return BatteryThermalMsg{Data: data}
	}
}

func parseBatteryFull(output string) BatteryThermal {
	var b BatteryThermal

	for _, line := range strings.Split(output, "\n") {
		key, val, ok := strings.Cut(strings.TrimSpace(line), ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)

		switch strings.TrimSpace(key) {
		case "level":
			b.Level = val + "%"
		case "status":
			b.Status = batteryStatusName(val)
		case "health":
			b.Health = batteryHealthName(val)
		case "voltage":
			b.Voltage = formatVoltage(val)
		case "temperature":
			b.Temperature = formatTenthsC(val)
		case "technology":
			b.Technology = val
		case "plugged":
			b.Plugged = batteryPluggedName(val)
		}
	}

	return b
}

func parseThermalZones(output string) []ThermalZone {
	var zones []ThermalZone

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		name, raw, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		raw = strings.TrimSpace(raw)
		if name == "" || raw == "" {
			continue
		}

		milli, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || milli <= 0 {
			continue
		}

		zones = append(zones, ThermalZone{
			Type: name,
			Temp: fmt.Sprintf("%.1f°C", float64(milli)/1000),
		})
		if len(zones) >= maxThermalZones {
			break
		}
	}

	return zones
}

func formatTenthsC(val string) string {
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.1f°C", float64(n)/10)
}

func formatVoltage(val string) string {
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil || n <= 0 {
		return ""
	}
	v := float64(n) / 1000 // millivolts
	if n > 100000 {
		v = float64(n) / 1e6 // some devices report microvolts
	}
	return fmt.Sprintf("%.2f V", v)
}

func batteryStatusName(code string) string {
	switch code {
	case "1":
		return "Unknown"
	case "2":
		return "Charging"
	case "3":
		return "Discharging"
	case "4":
		return "Not charging"
	case "5":
		return "Full"
	default:
		return code
	}
}

func batteryHealthName(code string) string {
	switch code {
	case "1":
		return "Unknown"
	case "2":
		return "Good"
	case "3":
		return "Overheat"
	case "4":
		return "Dead"
	case "5":
		return "Over voltage"
	case "6":
		return "Failure"
	case "7":
		return "Cold"
	default:
		return code
	}
}

func batteryPluggedName(code string) string {
	switch code {
	case "0":
		return "Unplugged"
	case "1":
		return "AC"
	case "2":
		return "USB"
	case "4":
		return "Wireless"
	default:
		return code
	}
}
