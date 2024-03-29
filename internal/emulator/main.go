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
	// registers are the CPU's 16-bit registers
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

func NewCPU8080(program *[]byte) *CPU8080 {
	vm := &CPU8080{
		Logger: log.New(os.Stdout),
	}
	// Put the program into memory at the location it wants to be
	copy(vm.memory[vm.Options.ProgramStartAddress:], *program)
	vm.programSize = len(*program) + int(vm.Options.ProgramStartAddress)

	// Define all supported opcodes
	vm.opcodeTable = map[byte]opcodeExec{
		0x00: vm.nop,
		0x01: vm.loadBC,
		0x05: vm.decB,
		0x06: vm.moveB,
		0x11: vm.loadDE,
		0x13: vm.incDE,
		0x1A: vm.loadAXD,
		0x21: vm.loadHL,
		0x23: vm.incHL,
		0x31: vm.loadSP,
		0x77: vm.storeHLA,
		0xC2: vm.jumpNZ,
		0xC3: vm.jump,
		// 0xC9: vm.ret,
		0xCD: vm.call,
	}

	return vm
}

func (vm *CPU8080) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	vm.nextOpCode()
	if int(vm.pc) == vm.programSize {
		return ebiten.Termination
	}

	return nil
}

func (vm *CPU8080) Draw(screen *ebiten.Image) {

}
func (vm *CPU8080) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
