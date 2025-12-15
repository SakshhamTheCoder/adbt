package state

import "adbt/internal/adb"

type AppState struct {
	SelectedDevice *adb.Device
	Devices        []adb.Device
	Width          int
	Height         int
}

func New() *AppState {
	return &AppState{
		Devices: []adb.Device{},
	}
}

func (s *AppState) SelectDevice(device *adb.Device) {
	s.SelectedDevice = device
}

func (s *AppState) HasDevice() bool {
	return s.SelectedDevice != nil
}

func (s *AppState) DeviceSerial() string {
	if s.SelectedDevice == nil {
		return ""
	}
	return s.SelectedDevice.Serial
}
