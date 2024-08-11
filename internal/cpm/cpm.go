// package cpm provides a dumb-down CP/M hardware environment to execute the TST8080 rom.
// This ROM is an excellent test rom that examines many 8080 codes for correctness.

package cpm

import (
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/TST8080.COM
var romFile embed.FS

type CPMHardware struct {
	rom []byte
}

func NewCPMHardware() *CPMHardware {
	romData, _ := romFile.ReadFile("assets/TST8080.COM")

	return &CPMHardware{
		rom: romData,
	}
}

func (cpm *CPMHardware) In(addr byte) (byte, error) {
	return 0, nil
}

func (cpm *CPMHardware) HandleSystemCall(vm *emulator.CPU8080) {
	if vm.PC == 0x0005 {
		switch vm.Registers.C {
		case 0x02:
			// print a single character
			fmt.Printf("%s", string(vm.Registers.E))
		case 0x09:
			// print a string
			start := (uint16(vm.Registers.D) << 8) | uint16(vm.Registers.E)
			end := start
			for {
				c := vm.Memory[end]
				if string(c) == "$" {
					break
				}
				fmt.Printf("%s", string(c))
				end++
			}
		}
		vm.PC++
	} else if vm.PC == 0 {
		os.Exit(0)
	}
}

func (cpm *CPMHardware) Out(addr byte, value byte) error {
	return nil
}

func (cpm *CPMHardware) InDeviceName(port byte) string {
	return ""
}

func (cpm *CPMHardware) OutDeviceName(port byte) string {
	return ""
}

func (cpm *CPMHardware) InterruptConditions() []emulator.Interrupt {
	return nil
}

func (cpm *CPMHardware) CyclesPerFrame() int {
	return 33334
}

func (cpm *CPMHardware) Draw(*ebiten.Image) {
}

func (cpm *CPMHardware) Init(memory *[65536]byte) {
	memory[0x0007] = 0xC9 // RET
}
func (cpm *CPMHardware) Width() int {
	return 224
}
func (cpm *CPMHardware) Height() int {
	return 256
}
func (cpm *CPMHardware) Scale() int {
	return 3
}
func (cpm *CPMHardware) StartAddress() int {
	return 0x100
}
func (cpm *CPMHardware) ROM() []byte {
	return cpm.rom
}
func (cpm *CPMHardware) FrameDuration() time.Duration {
	return 17 * time.Millisecond
}
func (cpm *CPMHardware) Cleanup() {
	//no-op
}
