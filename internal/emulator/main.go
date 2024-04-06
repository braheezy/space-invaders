package emulator

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

type CPU8080 struct {
	// pc is the current program counter, the address of the next instruction to be executed
	pc uint16
	// programData is a pointer to bytes containing the device memory (64kb)
	memory [64 * 1024]byte
	// programSize is the number of bytes in the program
	programSize int
	// registers are the CPU's 8-bit registers
	registers registers
	// sp is the stack pointer, the index to memory
	sp uint16
	// flags indicate the results of arithmetic and logical operations
	flags flags
	// Logger object to use
	Logger *log.Logger
	// Lookup table of opcode functions
	opcodeTable map[byte]opcodeExec
	// Options are the current options to use on the emulator
	Options CPU8080Options
	// For timing sync
	cycleCount int
	// io is the struct holding IO device interface methods
	IO HardwareIO
}

type CPU8080Options struct {
	ProgramStartAddress uint16
}

type registers struct {
	A, B, C, D, E, H, L byte
}

type flags struct {
	Z bool
	S bool
	H bool
	C bool
	P bool
}

// toByte packs flags according to the PSW layout to be pushed onto the stack.
func (f *flags) toByte() byte {
	var b byte
	// Set the Sign flag in the highest bit (bit 7)
	if f.S {
		b |= 1 << 7
	}
	// Set the Zero flag in bit 6
	if f.Z {
		b |= 1 << 6
	}
	// Set the Auxiliary Carry flag in bit 4
	if f.H {
		b |= 1 << 4
	}
	// Set the Parity flag in bit 2
	if f.P {
		b |= 1 << 2
	}
	// Bit 1 is always 1
	b |= 1 << 1
	// Set the Carry flag in bit 0
	if f.C {
		b |= 1
	}
	return b
}

func carrySub(value, subtrahend byte) bool {
	return value < subtrahend
}

func auxCarrySub(value, subtrahend byte) bool {
	// Check if borrow is needed from higher nibble to lower nibble
	return (value & 0xF) < (subtrahend & 0xF)
}
func auxCarryAdd(a, b byte) bool {
	// Check if carry is needed from higher nibble to lower nibble
	return (a & 0xF) > (b & 0xF)
}

func parity(x uint16) bool {
	y := x ^ (x >> 1)
	y = y ^ (y >> 2)
	y = y ^ (y >> 4)
	y = y ^ (y >> 8)

	// Rightmost bit of y holds the parity value
	// if (y&1) is 1 then parity is odd else even
	return y&1 > 0
}
func (fl *flags) setZ(value uint16) {
	fl.Z = value == 0
}
func (fl *flags) setS(value uint16) {
	fl.S = value&0x80 != 0
}
func (fl *flags) setP(value uint16) {
	fl.P = parity(value)
}

type opcodeExec func([]byte)

func NewCPU8080(program *[]byte, io HardwareIO) *CPU8080 {
	vm := &CPU8080{
		Logger: log.New(os.Stdout),
		IO:     io,
	}
	// Put the program into memory at the location it wants to be
	copy(vm.memory[vm.Options.ProgramStartAddress:], *program)
	vm.programSize = len(*program) + int(vm.Options.ProgramStartAddress)

	// Define all supported opcodes
	vm.opcodeTable = map[byte]opcodeExec{
		0x00: vm.nop,
		0x01: vm.load_BC,
		0x05: vm.dec_B,
		0x06: vm.moveI_B,
		0x09: vm.dad_BC,
		0x0D: vm.dec_C,
		0x0E: vm.moveI_C,
		0x0F: vm.rrc,
		0x11: vm.load_DE,
		0x13: vm.inc_DE,
		0x19: vm.dad_DE,
		0x1A: vm.load_DEA,
		0x21: vm.load_HL,
		0x23: vm.inc_HL,
		0x26: vm.moveI_H,
		0x29: vm.dad_HL,
		0x31: vm.load_SP,
		0x36: vm.moveI_HL,
		0x56: vm.moveHL_D,
		0x5E: vm.moveHL_E,
		0x66: vm.moveHL_H,
		0x6F: vm.move_AL,
		0x77: vm.load_HLA,
		0x7A: vm.move_DA,
		0x7C: vm.move_HA,
		0x7E: vm.moveHL_A,
		0xC2: vm.jump_NZ,
		0xC1: vm.pop_BC,
		0xC3: vm.jump,
		0xC5: vm.push_BC,
		0xC9: vm.ret,
		0xCD: vm.call,
		0xD1: vm.pop_DE,
		0xD3: vm.out,
		0xD5: vm.push_DE,
		0xE1: vm.pop_HL,
		0xE5: vm.push_HL,
		0xEB: vm.xchg,
		0xF5: vm.push_AF,
		0xFE: vm.cmp,
	}

	return vm
}

const (
	cyclesPerFrame = 33334 // Total cycles per frame, split into two halves
)

func (vm *CPU8080) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	vm.cycleCount = 0 // Reset cycle count
	vm.runCycles(cyclesPerFrame / 2)
	vm.performMidScreenInterrupt()

	vm.runCycles(cyclesPerFrame)
	vm.performFullScreenInterrupt()

	return nil
}

func (vm *CPU8080) Draw(screen *ebiten.Image) {

}
func (vm *CPU8080) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
