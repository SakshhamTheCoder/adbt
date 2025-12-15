package adb

import (
	"os/exec"
)

func StartScrcpy(serial string) error {
	cmd := exec.Command("scrcpy", "-s", serial)
	return cmd.Start()
}

func Reboot(serial string) error {
	_, err := ExecuteCommand(serial, "reboot")
	return err
}

func RebootRecovery(serial string) error {
	_, err := ExecuteCommand(serial, "reboot", "recovery")
	return err
}

func RebootBootloader(serial string) error {
	_, err := ExecuteCommand(serial, "reboot", "bootloader")
	return err
}
