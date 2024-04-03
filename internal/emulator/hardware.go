package emulator

import "fmt"

type HardwareIO interface {
	In(addr byte) byte
	Out(addr byte, data byte)
	DeviceName(port byte) string
}

type SpaceInvadersHardware struct {
	watchdogTimer byte
}

func (si *SpaceInvadersHardware) In(addr byte) byte {
	return 0
}

func (si *SpaceInvadersHardware) Out(addr byte, value byte) {
	if addr == 0x03 {
		si.watchdogTimer = value
	}
}
func (si *SpaceInvadersHardware) DeviceName(port byte) string {
	switch port {
	case 0x06:
		return "WATCHDOG"
	default:
		return fmt.Sprintf("$%02X", port) // Default to hex representation if unknown
	}
}

func NewSpaceInvadersHardware() *SpaceInvadersHardware {
	return &SpaceInvadersHardware{}
}
