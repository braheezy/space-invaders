package emulator

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

// stateCounts represents the number of clock cycles each 8080 CPU instruction takes to execute.
// The array index corresponds to the opcode, and the value at each index is the cycle count for that opcode.
// This is used in the main execution loop to track the total number of cycles run and ensure accurate timing.
var stateCounts = []int{
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 00..0f
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 00..1f
	4, 10, 16, 5, 5, 5, 7, 4, 4, 10, 16, 5, 5, 5, 7, 4, // 20..2f
	4, 10, 13, 5, 10, 10, 10, 4, 4, 10, 13, 5, 5, 5, 7, 4, // 30..3f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 40..4f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 50..5f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 60..6f
	7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 7, 5, // 70..7f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 80..8f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 90..9f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // a0..af
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // b0..bf
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 17, 7, 11, // c0..cf
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 17, 7, 11, // d0..df
	5, 10, 10, 18, 11, 11, 7, 11, 5, 5, 10, 5, 11, 17, 7, 11, // e0..ef
	5, 10, 10, 4, 11, 11, 7, 11, 5, 5, 10, 4, 11, 17, 7, 11, // f0..ff
}

// opcodeExec is a function to execute for the current opcode
type opcodeExec func([]byte)

// CPU8080 emulates an Intel 8080 CPU processor
type CPU8080 struct {
	// PC is the current program counter, the address of the next instruction to be executed
	PC uint16
	// programData is a pointer to bytes containing the device Memory (64kb)
	Memory [64 * 1024]byte
	// programSize is the number of bytes in the program
	programSize int
	// Registers are the CPU's 8-bit Registers
	Registers Registers
	// sp is the stack pointer, the index to memory
	sp uint16
	// flags indicate the results of arithmetic and logical operations
	flags flags
	// Logger object to use
	Logger *log.Logger
	// Lookup table of opcode functions
	opcodeTable map[byte]opcodeExec
	// Options are the current options to use on the emulator
	Options EmulatorOptions
	// For timing sync
	cycleCount  int
	totalCycles int
	// Hardware is the struct holding HardwareIO device interface methods
	Hardware HardwareIO
	// Mutex for thread-safe access to the CPU state
	mu sync.Mutex
	// Whether or not interrupts are currently being handled
	interruptsEnabled bool
	// InterruptRequest is a channel to request an interrupt by sending an opcode.
	InterruptRequest chan byte
}

// EmulatorOptions describe tunable settings about emulator execution
type EmulatorOptions struct {
	UnlimitedTPS bool
}

// Registers are the 7 primary registers for the 8080.
type Registers struct {
	A, B, C, D, E, H, L byte
}

// flags represent the conditions bits set after data operations. These can be checked by later instructions to
// affect execution.
type flags struct {
	// Zero flag set if result is zero
	Z bool
	// Sign flag set if result is negative
	S bool
	// Auxillary Carry flag set if a carry out of bit 3 occurs
	H bool
	// Carry flag set during arithmetic operations
	// Addition: if sum exceeds max byte value, carry occurs
	// Subtraction: if result is less than zero, carry occurs
	C bool
	// Parity flag set if number of bits in result is even
	P bool
}

// NewEmulator creates a new emulator, combing the provided HardwareIO with a CPU8080
func NewEmulator(io HardwareIO) *CPU8080 {
	// Load the ROM from the hardware
	program := io.ROM()

	// Initialize emulator virtual machine
	vm := &CPU8080{
		Logger:            log.New(os.Stdout),
		Hardware:          io,
		InterruptRequest:  make(chan byte, 1),
		interruptsEnabled: true,
	}
	start := io.StartAddress()
	// Put the program into memory at the location it wants to be
	copy(vm.Memory[start:], program)

	// Calculate program size for graceful termination
	vm.programSize = len(program) + start
	// Initialize program counter to start address
	vm.PC = uint16(start)

	// Give the hardware initialization time with the hardware
	vm.Hardware.Init(&vm.Memory)

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
		0x0B: vm.dcx_B,
		0x0C: vm.inr_C,
		0x0D: vm.dcr_C,
		0x0E: vm.moveImm_C,
		0x0F: vm.rrc,
		0x11: vm.load_DE,
		0x12: vm.stax_D,
		0x14: vm.inr_D,
		0x15: vm.dcr_D,
		0x13: vm.inx_D,
		0x16: vm.moveImm_D,
		0x17: vm.ral,
		0x19: vm.dad_D,
		0x1A: vm.loadAddr_D,
		0x1B: vm.dcx_D,
		0x1C: vm.inr_E,
		0x1D: vm.dcr_E,
		0x1E: vm.moveImm_E,
		0x1F: vm.rar,
		0x21: vm.load_HL,
		0x22: vm.store_HL,
		0x23: vm.inx_H,
		0x24: vm.inr_H,
		0x25: vm.dcr_H,
		0x26: vm.moveImm_H,
		0x27: vm.daa,
		0x29: vm.dad_H,
		0x2A: vm.loadImm_HL,
		0x2C: vm.inr_L,
		0x2B: vm.dcx_H,
		0x2D: vm.dcr_L,
		0x2E: vm.moveImm_L,
		0x2F: vm.cma,
		0x31: vm.load_SP,
		0x32: vm.store_A,
		0x33: vm.inx_SP,
		0x34: vm.inr_M,
		0x35: vm.dcr_M,
		0x36: vm.moveImm_M,
		0x37: vm.set_C,
		0x39: vm.dad_SP,
		0x3A: vm.load_A,
		0x3B: vm.dcx_SP,
		0x3C: vm.inr_A,
		0x3D: vm.dcr_A,
		0x3E: vm.moveImm_A,
		0x3F: vm.cmc,
		0x41: vm.move_BC,
		0x42: vm.move_BD,
		0x43: vm.move_BE,
		0x44: vm.move_BH,
		0x45: vm.move_BL,
		0x46: vm.move_BM,
		0x47: vm.move_BA,
		0x48: vm.move_CB,
		0x4A: vm.move_CD,
		0x4B: vm.move_CE,
		0x4C: vm.move_CH,
		0x4D: vm.move_CL,
		0x4E: vm.move_CM,
		0x4F: vm.move_CA,
		0x50: vm.move_DB,
		0x51: vm.move_DC,
		0x53: vm.move_DE,
		0x54: vm.move_DH,
		0x55: vm.move_DL,
		0x56: vm.move_DM,
		0x57: vm.move_DA,
		0x58: vm.move_EB,
		0x59: vm.move_EC,
		0x5A: vm.move_ED,
		0x5C: vm.move_EH,
		0x5D: vm.move_EL,
		0x5E: vm.move_EM,
		0x5F: vm.move_EA,
		0x60: vm.move_HB,
		0x61: vm.move_HC,
		0x62: vm.move_HD,
		0x63: vm.move_HE,
		0x65: vm.move_HL,
		0x66: vm.move_HM,
		0x67: vm.move_HA,
		0x68: vm.move_LB,
		0x69: vm.move_LC,
		0x6A: vm.move_LD,
		0x6B: vm.move_LE,
		0x6C: vm.move_LH,
		0x6E: vm.move_LM,
		0x6F: vm.move_LA,
		0x70: vm.move_MB,
		0x71: vm.move_MC,
		0x72: vm.move_MD,
		0x73: vm.move_ME,
		0x74: vm.move_MH,
		0x75: vm.move_ML,
		0x77: vm.move_MA,
		0x78: vm.move_AB,
		0x79: vm.move_AC,
		0x7A: vm.move_AD,
		0x7B: vm.move_AE,
		0x7C: vm.move_AH,
		0x7D: vm.move_AL,
		0x7E: vm.move_AM,
		0x80: vm.add_B,
		0x81: vm.add_C,
		0x82: vm.add_D,
		0x83: vm.add_E,
		0x84: vm.add_H,
		0x85: vm.add_L,
		0x86: vm.add_M,
		0x87: vm.add_A,
		0x88: vm.adc_B,
		0x89: vm.adc_C,
		0x8A: vm.adc_D,
		0x8B: vm.adc_E,
		0x8C: vm.adc_H,
		0x8D: vm.adc_L,
		0x8E: vm.adc_M,
		0x8F: vm.adc_A,
		0x90: vm.sub_B,
		0x91: vm.sub_C,
		0x92: vm.sub_D,
		0x93: vm.sub_E,
		0x94: vm.sub_H,
		0x95: vm.sub_L,
		0x96: vm.sub_M,
		0x97: vm.sub_A,
		0x98: vm.sbb_B,
		0x99: vm.sbb_C,
		0x9A: vm.sbb_D,
		0x9B: vm.sbb_E,
		0x9C: vm.sbb_H,
		0x9D: vm.sbb_L,
		0x9E: vm.sbb_M,
		0x9F: vm.sbb_A,
		0xA0: vm.ana_B,
		0xA1: vm.ana_C,
		0xA2: vm.ana_D,
		0xA3: vm.ana_E,
		0xA4: vm.ana_H,
		0xA5: vm.ana_L,
		0xA6: vm.ana_M,
		0xA7: vm.ana_A,
		0xA8: vm.xra_B,
		0xA9: vm.xra_C,
		0xAA: vm.xra_D,
		0xAB: vm.xra_E,
		0xAC: vm.xra_H,
		0xAD: vm.xra_L,
		0xAE: vm.xra_M,
		0xAF: vm.xra_A,
		0xB0: vm.ora_B,
		0xB1: vm.ora_C,
		0xB2: vm.ora_D,
		0xB3: vm.ora_E,
		0xB4: vm.ora_H,
		0xB5: vm.ora_L,
		0xB6: vm.ora_M,
		0xB7: vm.ora_A,
		0xB8: vm.cmp_B,
		0xB9: vm.cmp_C,
		0xBA: vm.cmp_D,
		0xBB: vm.cmp_E,
		0xBC: vm.cmp_H,
		0xBD: vm.cmp_L,
		0xBE: vm.cmp_M,
		0xBF: vm.cmp_A,
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
		0xCE: vm.aci,
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
		0xDC: vm.call_C,
		0xDE: vm.sbi,
		0xE0: vm.ret_PO,
		0xE1: vm.pop_HL,
		0xE2: vm.jump_PO,
		0xE3: vm.xthl,
		0xE4: vm.call_PO,
		0xE5: vm.push_HL,
		0xE9: vm.pchl,
		0xEA: vm.jump_PE,
		0xEB: vm.xchg,
		0xE6: vm.and,
		0xE8: vm.ret_PE,
		0xEC: vm.call_PE,
		0xEE: vm.xri,
		0xF1: vm.pop_AF,
		0xF2: vm.jump_P,
		0xF3: vm.di,
		0xFA: vm.jump_M,
		0xFB: vm.ei,
		0xF0: vm.ret_P,
		0xF4: vm.call_P,
		0xF5: vm.push_AF,
		0xF6: vm.ori,
		0xF8: vm.ret_M,
		0xF9: vm.sphl,
		0xFC: vm.call_M,
		0xFE: vm.cmp,
	}

	return vm
}

// StartInterruptRoutines starts a goroutine for each interrupt condition.
// These goroutines will check if the interrupt condition is met and request an interrupt if so.
// While the goroutines run in an infinite loop, the condition is only checked at the specified cycle interval.
func (vm *CPU8080) StartInterruptRoutines() {
	for _, condition := range vm.Hardware.InterruptConditions() {
		go func(cond Interrupt) {
			ticker := time.NewTicker(time.Duration(cond.Cycle) * time.Nanosecond)
			for {
				<-ticker.C
				if vm.interruptsEnabled && vm.cycleCount >= cond.Cycle {
					// Run the interrupt routine
					cond.Action(vm)
				}
			}
		}(condition)
	}
}

// runCycles executes the CPU for cycleCount amount of times.
// This is the main execution loop of the emulator.
func (vm *CPU8080) runCycles(cycleCount int) {
	// Record when the frame started, in case we need to slow down later
	// for historically accurate speeds.
	var startTime time.Time
	if !vm.Options.UnlimitedTPS {
		startTime = time.Now()
	}

	for vm.cycleCount < cycleCount {
		select {
		// After every opcode execution, check if an interrupt was requested
		case opcode := <-vm.InterruptRequest:
			vm.handleInterrupt(opcode)
		// Process next opcode
		default:
			if int(vm.PC) >= vm.programSize {
				// There's nothing left to process!
				break
			}

			// Some hardware perform IO operations through system calls instead
			// of IN/OUT opcodes. Allow that to happen here.
			vm.Hardware.HandleSystemCall(vm)

			// Parse the next 3 bytes for this opcode execution.
			currentCode := vm.Memory[vm.PC : vm.PC+3]
			op := currentCode[0]
			vm.PC++
			vm.cycleCount += stateCounts[op]
			vm.totalCycles += stateCounts[op]

			if opcodeFunc, exists := vm.opcodeTable[op]; exists {
				opcodeFunc(currentCode[1:])
			} else {
				vm.Logger.Fatal("unsupported", "address", fmt.Sprintf("%04X", vm.PC-1), "opcode", fmt.Sprintf("%02X", op), "totalCycles", vm.totalCycles)
			}
		}
	}

	// Handle slowdown for accurate speed emulation
	if !vm.Options.UnlimitedTPS {
		elapsed := time.Since(startTime)
		if remaining := vm.Hardware.FrameDuration() - elapsed; remaining > 0 {
			time.Sleep(remaining)
		}
	}
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

// fromByte unpacks flags from the stack
func fromByte(b byte) *flags {
	return &flags{
		S: b&(1<<7) != 0,
		Z: b&(1<<6) != 0,
		H: b&(1<<4) != 0,
		P: b&(1<<2) != 0,
		C: b&1 != 0,
	}
}

// carrySub returns true if a carry would happen if subtrahend is subtracted from value.
func carrySub(value, subtrahend byte) bool {
	return value < subtrahend
}

// carryAdd returns true if a carry would happen if addend is added to value.
func carryAdd(value, addend byte) bool {
	return uint16(value)+uint16(addend) > 0xFF
}

// auxCarrySub returns true if auxillary carry would happen if subtrahend is subtracted from value.
func auxCarrySub(value, subtrahend byte) bool {
	// Check if borrow is needed from higher nibble to lower nibble
	return (value & 0xF) < (subtrahend & 0xF)
}

// auxCarryAdd returns true if auxillary carry would happen if addend is added to value.
func auxCarryAdd(value, addend byte) bool {
	// Check if carry is needed from higher nibble to lower nibble
	return (value&0xF)+(addend&0xF) > 0xF
}

// parity returns true if the number of bits in x is even.
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
func toUint16(high, low byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// Update fulfills the Game interface for ebiten.
// This runs the emulator for one frame.
func (vm *CPU8080) Update() error {

	// Reset cycle count
	vm.cycleCount = 0
	// Execute opcodes
	vm.runCycles(vm.Hardware.CyclesPerFrame())

	return nil
}

// Draw fulfills the Game interface for ebiten
func (vm *CPU8080) Draw(screen *ebiten.Image) {
	// Use hardware to draw on the display
	vm.Hardware.Draw(screen)
}

// y fulfills the Game interface for ebiten
func (vm *CPU8080) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
