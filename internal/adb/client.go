package adb

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// commandTimeout bounds one-shot adb calls so an offline or unauthorized device
// can't hang the UI forever. Streaming commands (logcat, scrcpy) don't go
// through ExecuteCommand and are unaffected.
const commandTimeout = 30 * time.Second

func ExecuteCommand(serial string, args ...string) ([]byte, error) {
	var cmdArgs []string
	if serial != "" {
		cmdArgs = append(cmdArgs, "-s", serial)
	}
	cmdArgs = append(cmdArgs, args...)

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "adb", cmdArgs...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return out, fmt.Errorf("adb timed out after %s (is the device responsive?)", commandTimeout)
	}
	if err != nil {
		return out, fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
	}
	return out, nil
}

// EnsureAvailable reports whether the adb binary can be found on PATH.
func EnsureAvailable() error {
	if _, err := exec.LookPath("adb"); err != nil {
		return fmt.Errorf("adb not found in PATH")
	}
	return nil
}

func GetProperty(serial, prop string) (string, error) {
	out, err := ExecuteCommand(serial, "shell", "getprop", prop)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func ParseLines(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
