package emulator

import (
	"os"
	"sync"
	"time"

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
	cycleCount  int
	totalCycles int
	// io is the struct holding Hardware device interface methods
	Hardware          HardwareIO
	mu                sync.Mutex // Mutex for thread-safe access to the CPU state
	InterruptsEnabled bool
	InterruptRequest  chan byte
}

type CPU8080Options struct {
	UnlimitedTPS bool
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

func fromByte(b byte) *flags {
	return &flags{
		S: b&(1<<7) != 0, // Check if the Sign flag bit (bit 7) is set
		Z: b&(1<<6) != 0, // Check if the Zero flag bit (bit 6) is set
		H: b&(1<<4) != 0, // Check if the Auxiliary Carry flag bit (bit 4) is set
		P: b&(1<<2) != 0, // Check if the Parity flag bit (bit 2) is set
		C: b&1 != 0,      // Check if the Carry flag bit (bit 0) is set
	}
}

func carrySub(value, subtrahend byte) bool {
	return value < subtrahend
}

func carryAdd(value, addend byte) bool {
	return uint16(value)+uint16(addend) > 0xFF
}

func auxCarrySub(value, subtrahend byte) bool {
	// Check if borrow is needed from higher nibble to lower nibble
	return (value & 0xF) < (subtrahend & 0xF)
}
func auxCarryAdd(value, addend byte) bool {
	// Check if carry is needed from higher nibble to lower nibble
	return (value&0xF)+(addend&0xF) > 0xF
}

func parity(x uint16) bool {
	y := x ^ (x >> 1)
	y = y ^ (y >> 2)
	y = y ^ (y >> 4)
	y = y ^ (y >> 8)

	// Rightmost bit of y holds the parity value
	// if (y&1) is 1 then parity is odd else even
	return y&1 == 0
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

func (vm *CPU8080) StartInterruptRoutines() {
	for _, condition := range vm.Hardware.InterruptConditions() {
		go func(cond InterruptCondition) {
			ticker := time.NewTicker(time.Duration(cond.Cycle) * time.Nanosecond)
			for {
				<-ticker.C
				if vm.InterruptsEnabled && vm.cycleCount >= cond.Cycle {
					cond.Action(vm)
				}
			}
		}(condition)
	}
}

type opcodeExec func([]byte)

func NewCPU8080(program *[]byte, io HardwareIO) *CPU8080 {
	vm := &CPU8080{
		Logger:            log.New(os.Stdout),
		Hardware:          io,
		InterruptRequest:  make(chan byte, 1),
		InterruptsEnabled: true,
	}
	// Put the program into memory at the location it wants to be
	copy(vm.memory[io.StartAddress():], *program)
	vm.programSize = len(*program) + int(io.StartAddress())

	if vm.Hardware != nil {
		vm.Hardware.Init(&vm.memory)
	}

	// Define all supported opcodes
	vm.opcodeTable = map[byte]opcodeExec{
		0x00: vm.nop,
		0x01: vm.load_BC,
		0x02: vm.stax_B,
		0x03: vm.inx_B,
		0x04: vm.inr_B,
		0x05: vm.dcr_B,
		0x06: vm.moveImm_B,
		0x07: vm.rlc,
		0x09: vm.dad_B,
		0x0A: vm.loadAddr_B,
		0x0C: vm.inr_C,
		0x0D: vm.dcr_C,
		0x0E: vm.moveImm_C,
		0x0F: vm.rrc,
		0x11: vm.load_DE,
		0x14: vm.inr_D,
		0x15: vm.dcr_D,
		0x13: vm.inx_D,
		0x16: vm.moveImm_D,
		0x19: vm.dad_D,
		0x1A: vm.loadAddr_D,
		0x1F: vm.rar,
		0x21: vm.load_HL,
		0x22: vm.store_HL,
		0x23: vm.inx_H,
		0x26: vm.moveImm_H,
		0x27: vm.daa,
		0x29: vm.dad_H,
		0x2A: vm.loadImm_HL,
		0x2B: vm.dcx_H,
		0x2E: vm.moveImm_L,
		0x31: vm.load_SP,
		0x32: vm.store_A,
		0x35: vm.dcr_M,
		0x36: vm.moveImm_M,
		0x37: vm.set_C,
		0x3A: vm.load_A,
		0x3C: vm.inr_A,
		0x3D: vm.dcr_A,
		0x3E: vm.moveImm_A,
		0x46: vm.move_BM,
		0x47: vm.move_BA,
		0x4E: vm.move_CM,
		0x4F: vm.move_CA,
		0x56: vm.move_DM,
		0x57: vm.move_DA,
		0x5E: vm.move_EM,
		0x5F: vm.move_EA,
		0x61: vm.move_HC,
		0x66: vm.move_HM,
		0x67: vm.move_HA,
		0x68: vm.move_LB,
		0x6F: vm.move_LA,
		0x70: vm.move_MH,
		0x77: vm.move_MA,
		0x78: vm.move_AB,
		0x79: vm.move_AC,
		0x7A: vm.move_AD,
		0x7B: vm.move_AE,
		0x7C: vm.move_AH,
		0x7D: vm.move_AL,
		0x7E: vm.move_AM,
		0x80: vm.add_B,
		0x83: vm.add_E,
		0x85: vm.add_L,
		0x86: vm.add_M,
		0xA0: vm.ana_B,
		0xA7: vm.ana_A,
		0xAF: vm.xra_A,
		0xB0: vm.ora_B,
		0xB4: vm.ora_H,
		0xB6: vm.ora_M,
		0xB8: vm.cmp_B,
		0xBE: vm.cmp_M,
		0xC0: vm.ret_NZ,
		0xC2: vm.jump_NZ,
		0xC1: vm.pop_BC,
		0xC3: vm.jump,
		0xC4: vm.call_NZ,
		0xC5: vm.push_BC,
		0xC6: vm.adi,
		0xC8: vm.ret_Z,
		0xC9: vm.ret,
		0xCA: vm.jump_Z,
		0xCC: vm.call_Z,
		0xCD: vm.call,
		0xD0: vm.ret_NC,
		0xD1: vm.pop_DE,
		0xD2: vm.jump_NC,
		0xD3: vm.out,
		0xD4: vm.call_NC,
		0xD5: vm.push_DE,
		0xD6: vm.sui,
		0xD8: vm.ret_C,
		0xDA: vm.jump_C,
		0xDB: vm.in,
		0xDE: vm.sbi,
		0xE1: vm.pop_HL,
		0xE3: vm.xthl,
		0xE5: vm.push_HL,
		0xE9: vm.pchl,
		0xEB: vm.xchg,
		0xE6: vm.and,
		0xF1: vm.pop_AF,
		0xF3: vm.di,
		0xFA: vm.jump_m,
		0xFB: vm.ei,
		0xF4: vm.call_P,
		0xF5: vm.push_AF,
		0xF6: vm.ori,
		0xFE: vm.cmp,
	}

	return vm
}

func (vm *CPU8080) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	vm.cycleCount = 0
	vm.runCycles(vm.Hardware.CyclesPerFrame())

	return nil
}

func (vm *CPU8080) Draw(screen *ebiten.Image) {
	vm.Hardware.Draw(screen)
}
func (vm *CPU8080) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
