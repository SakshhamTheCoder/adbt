package adb

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteCommand(serial string, args ...string) ([]byte, error) {
	var cmdArgs []string
	if serial != "" {
		cmdArgs = append(cmdArgs, "-s", serial)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("adb", cmdArgs...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return out, fmt.Errorf("%w: %s", err, bytes.TrimSpace(out))
	}
	return out, nil
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
