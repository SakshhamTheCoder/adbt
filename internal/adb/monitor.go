package adb

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// SystemStats holds raw counters from /proc files.
// The UI layer is responsible for calculating rates/percentages by comparing with previous samples.
type SystemStats struct {
	// CPU counters (jiffies)
	CPUTotal uint64
	CPUIdle  uint64

	// Memory info (kB)
	MemTotal     uint64
	MemUsed      uint64 // Total - Available
	MemAvailable uint64 // MemFree + Buffers + Cached (approx)

	// Network counters (bytes) - usually for wlan0 or total
	NetRxBytes uint64
	NetTxBytes uint64
}

type SystemStatsMsg struct {
	Stats SystemStats
	Error error
}

func GetSystemStatsCmd(serial string) tea.Cmd {
	return func() tea.Msg {
		var stats SystemStats
		var err error

		// 1. CPU
		out, err := ExecuteCommand(serial, "shell", "cat", "/proc/stat")
		if err == nil {
			localTotal, localIdle := parseCPUStats(string(out))
			stats.CPUTotal = localTotal
			stats.CPUIdle = localIdle
		}

		// 2. Memory
		out, err = ExecuteCommand(serial, "shell", "cat", "/proc/meminfo")
		if err == nil {
			t, a := parseMemInfo(string(out))
			stats.MemTotal = t
			stats.MemAvailable = a
			if t > a {
				stats.MemUsed = t - a
			}
		}

		// 3. Network
		out, err = ExecuteCommand(serial, "shell", "cat", "/proc/net/dev")
		if err == nil {
			rx, tx := parseNetDev(string(out))
			stats.NetRxBytes = rx
			stats.NetTxBytes = tx
		}

		return SystemStatsMsg{Stats: stats}
	}
}

func parseCPUStats(output string) (total, idle uint64) {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return
	}

	// first line: cpu  2255 34 2290 22625563 6290 127 456 0 0 0
	line := lines[0]
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return
	}

	// fields[1] = user, [2] = nice, [3] = system, [4] = idle, [5] = iowait, ...
	var values []uint64
	for _, f := range fields[1:] {
		v, _ := strconv.ParseUint(f, 10, 64)
		values = append(values, v)
	}

	var sum uint64
	for _, v := range values {
		sum += v
	}
	total = sum

	if len(values) >= 4 {
		idleVal := values[3] // idle
		if len(values) >= 5 {
			idleVal += values[4] // + iowait
		}
		idle = idleVal
	}

	return
}

func parseMemInfo(output string) (total, available uint64) {
	lines := strings.Split(output, "\n")
	var memFree, buffers, cached uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		val, _ := strconv.ParseUint(fields[1], 10, 64)

		switch key {
		case "MemTotal":
			total = val
		case "MemAvailable":
			available = val
		case "MemFree":
			memFree = val
		case "Buffers":
			buffers = val
		case "Cached":
			cached = val
		}
	}

	if available == 0 {
		available = memFree + buffers + cached
	}

	return
}

func parseNetDev(output string) (rx, tx uint64) {
	lines := strings.Split(output, "\n")
	// Inter-face   |   Receive ... | Transmit ...
	// wlan0: 123 456 ... or
	//  wlan0:123 ... (no space after colon)

	for _, line := range lines {
		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		iface := strings.TrimSpace(parts[0])
		if iface == "lo" {
			continue
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 9 { // Receive: bytes(0), packets(1), errs(2), drop(3), fifo(4), frame(5), compressed(6), multicast(7) | Transmit: bytes(8)
			continue
		}

		r, _ := strconv.ParseUint(fields[0], 10, 64)
		t, _ := strconv.ParseUint(fields[8], 10, 64)

		rx += r
		tx += t
	}
	return
}
