package state

import "github.com/SakshhamTheCoder/adbt/internal/adb"

type AppState struct {
	SelectedDeviceSerial string
	Devices              []adb.Device
	Width                int
	Height               int
}

func New() *AppState {
	return &AppState{
		Devices: []adb.Device{},
	}
}

func (s *AppState) SelectDevice(serial string) {
	s.SelectedDeviceSerial = serial
}

func (s *AppState) SelectedDevice() *adb.Device {
	if s.SelectedDeviceSerial == "" {
		return nil
	}

	for i := range s.Devices {
		if s.Devices[i].Serial == s.SelectedDeviceSerial {
			return &s.Devices[i]
		}
	}

	return nil
}

func (s *AppState) HasDevice() bool {
	return s.SelectedDevice() != nil
}

func (s *AppState) DeviceSerial() string {
	if device := s.SelectedDevice(); device != nil {
		return device.Serial
	}
	return ""
}
