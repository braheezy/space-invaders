package emulator

import (
	"fmt"
)

type InterruptCondition struct {
	Cycle  int
	Action func(*CPU8080)
}

type HardwareIO interface {
	In(addr byte) byte
	Out(addr byte, data byte)
	DeviceName(port byte) string
	InterruptConditions() []InterruptCondition
	CyclesPerFrame() int
}

type SpaceInvadersHardware struct {
	watchdogTimer  byte
	cyclesPerFrame int
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
func (si *SpaceInvadersHardware) InterruptConditions() []InterruptCondition {
	return []InterruptCondition{
		{
			// Mid screen interrupt
			Cycle: si.cyclesPerFrame / 2,
			Action: func(vm *CPU8080) {
				// RST 8
				vm.interruptRequest <- 0xCF
			},
		},
		{
			// VBLANK interrupt
			Cycle: si.cyclesPerFrame,
			Action: func(vm *CPU8080) {
				// RST 10
				vm.interruptRequest <- 0xD7
			},
		},
	}
}

func (si *SpaceInvadersHardware) CyclesPerFrame() int {
	return si.cyclesPerFrame // The constant value for Space Invaders
}
func NewSpaceInvadersHardware() *SpaceInvadersHardware {

	return &SpaceInvadersHardware{
		cyclesPerFrame: 33334,
	}
}
