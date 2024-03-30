package emulator

import (
	"time"
)

func (vm *CPU8080) runCycles(cycleCount int) {
	startTime := time.Now()

	for vm.cycleCount < cycleCount {
		if int(vm.pc) >= vm.programSize {
			break
		}
		currentCode := vm.memory[vm.pc : vm.pc+3]

		op := currentCode[0]
		vm.pc++
		vm.cycleCount++

		if opcodeFunc, exists := vm.opcodeTable[op]; exists {
			opcodeFunc(currentCode[1:])
		} else {
			vm.Logger.Fatalf("Unsupported opcode: %02X", op)
		}
	}

	elapsed := time.Since(startTime)
	if remaining := (17 * time.Millisecond) - elapsed; remaining > 0 {
		time.Sleep(remaining)
	}

}

func toUint16(code *[]byte) uint16 {
	return uint16((*code)[1])<<8 | uint16((*code)[0])
}

func auxCarrySub(a, b byte) bool {
	// Check if borrow is needed from higher nibble to lower nibble
	return (a & 0xF) < (b & 0xF)
}

// func auxCarryAdd(a, b byte) bool {
// 	// Check if carry is needed from higher nibble to lower nibble
// 	return (a & 0xF) > (b & 0xF)
// }

func parity(x uint16) bool {
	y := x ^ (x >> 1)
	y = y ^ (y >> 2)
	y = y ^ (y >> 4)
	y = y ^ (y >> 8)

	// Rightmost bit of y holds the parity value
	// if (y&1) is 1 then parity is odd else even
	return y&1 > 0
}
func (vm *CPU8080) performMidScreenInterrupt() {
	// Implement mid-screen interrupt tasks here
}

func (vm *CPU8080) performFullScreenInterrupt() {
	// Implement full-screen interrupt tasks here
}

// NOP: No operation.
func (vm *CPU8080) nop(data []byte) {
	vm.Logger.Debugf("[00] NOP")
}

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[C3] JMP \t$%04X", operand)
	vm.pc = operand
	vm.cycleCount += 2
}

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) loadSP(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[31] LD  \tSP,$%04X", operand)
	vm.sp = operand
	vm.pc += 2
	vm.cycleCount += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) loadBC(data []byte) {
	vm.Logger.Debugf("[01] LD  \tB,$%04X", toUint16(&data))
	vm.registers.C = data[0]
	vm.registers.B = data[1]
	vm.pc += 2
	vm.cycleCount += 2
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveB(data []byte) {
	vm.Logger.Debugf("[06] LD  \tB,$%02X", data[0])
	vm.registers.B = data[0]
	vm.pc++
	vm.cycleCount++
}

// CALL addr: Call subroutine at address
func (vm *CPU8080) call(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[CD] CALL\t$%04X", operand)
	vm.pc = operand
	vm.memory[vm.sp-1] = data[1]
	vm.memory[vm.sp-2] = data[0]
	vm.sp -= 2
	vm.cycleCount += 5
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) loadDE(data []byte) {
	vm.Logger.Debugf("[11] LD  \tDE,$%04X", toUint16(&data))
	vm.registers.E = data[0]
	vm.registers.D = data[1]
	vm.pc += 2
	vm.cycleCount += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) loadHL(data []byte) {
	vm.Logger.Debugf("[21] LD  \tHL,$%04X", toUint16(&data))
	vm.registers.L = data[0]
	vm.registers.H = data[1]
	vm.pc += 2
	vm.cycleCount += 2
}

// LDAX D: Load value from address in register pair D into accumulator.
func (vm *CPU8080) loadAXD(data []byte) {
	address := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	vm.Logger.Debugf("[1A] LD  \tA,(DE)")
	vm.registers.A = vm.memory[address]
	vm.cycleCount += 2
}

// MOV M, A: Move value from accumulator into register pair H.
func (vm *CPU8080) storeHLA(data []byte) {
	address := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	vm.Logger.Debugf("[77] LD  \t(HL),A ($%04X)", address)
	vm.memory[address] = vm.registers.A
	vm.cycleCount += 5
}

// INC H: Increment register pair H.
func (vm *CPU8080) inxHL(data []byte) {
	vm.Logger.Debugf("[23] INX \tHL")
	hl := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	hl++
	vm.registers.H = byte(hl >> 8)
	vm.registers.L = byte(hl & 0xFF)
	vm.cycleCount++
}

// INC D: Increment register pair D.
func (vm *CPU8080) inxDE(data []byte) {
	vm.Logger.Debugf("[23] INX \tDE")
	de := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	de++
	vm.registers.H = byte(de >> 8)
	vm.registers.L = byte(de & 0xFF)
	vm.cycleCount++
}

// DCR B: Decrement register B.
func (vm *CPU8080) decB(data []byte) {
	vm.Logger.Debugf("[05] DEC \tB")
	result := uint16(vm.registers.B) - 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.H = auxCarrySub(vm.registers.B, 1)
	vm.flags.setP(result)

	vm.registers.B--
	vm.cycleCount++
}

// JNZ addr: Jump if not zero.
func (vm *CPU8080) jumpNZ(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", operand)
	if !vm.flags.Z {
		vm.pc = operand
	} else {
		vm.pc += 2
	}
	vm.cycleCount += 2
}

// RET: Return from subroutine.
func (vm *CPU8080) ret(data []byte) {
	vm.Logger.Debugf("[C9] RET")
	vm.pc = toUint16(&[]byte{vm.memory[vm.sp+1], vm.memory[vm.sp]})
	vm.sp += 2
}
