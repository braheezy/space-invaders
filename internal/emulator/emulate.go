package emulator

import (
	"time"
)

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

func (vm *CPU8080) runCycles(cycleCount int) {
	startTime := time.Now()

	for vm.cycleCount < cycleCount {
		if int(vm.pc) >= vm.programSize {
			break
		}
		currentCode := vm.memory[vm.pc : vm.pc+3]

		op := currentCode[0]
		vm.pc++
		vm.cycleCount += stateCounts[op]

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
}

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) load_SP(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[31] LD  \tSP,$%04X", operand)
	vm.sp = operand
	vm.pc += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) load_BC(data []byte) {
	vm.Logger.Debugf("[01] LD  \tB,$%04X", toUint16(&data))
	vm.registers.C = data[0]
	vm.registers.B = data[1]
	vm.pc += 2
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveI_B(data []byte) {
	vm.Logger.Debugf("[06] LD  \tB,$%02X", data[0])
	vm.registers.B = data[0]
	vm.pc++
}

// MVI C, D8: Move 8-bit immediate value into register C.
func (vm *CPU8080) moveI_C(data []byte) {
	vm.Logger.Debugf("[0E] LD  \tC,$%02X", data[0])
	vm.registers.C = data[0]
	vm.pc++
}

// MVI H, D8: Move 8-bit immediate value into register H.
func (vm *CPU8080) moveI_H(data []byte) {
	vm.Logger.Debugf("[26] LD  \tH,$%02X", data[0])
	vm.registers.H = data[0]
	vm.pc++
}

// CALL addr: Call subroutine at address
func (vm *CPU8080) call(data []byte) {
	jumpAddress := toUint16(&data)
	returnAddress := vm.pc + 2
	vm.Logger.Debugf("[CD] CALL\t$%04X", jumpAddress)
	vm.memory[vm.sp-1] = byte(returnAddress >> 8)
	vm.memory[vm.sp-2] = byte(returnAddress & 0xFF)
	vm.pc = jumpAddress
	vm.sp -= 2
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) load_DE(data []byte) {
	vm.Logger.Debugf("[11] LD  \tDE,$%04X", toUint16(&data))
	vm.registers.E = data[0]
	vm.registers.D = data[1]
	vm.pc += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) load_HL(data []byte) {
	vm.Logger.Debugf("[21] LD  \tHL,$%04X", toUint16(&data))
	vm.registers.L = data[0]
	vm.registers.H = data[1]
	vm.pc += 2
}

// LDAX D: Load value from address in register pair D into accumulator.
func (vm *CPU8080) load_DEA(data []byte) {
	address := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	vm.Logger.Debugf("[1A] LD  \tA,(DE)")
	vm.registers.A = vm.memory[address]
}

// MOV M, A: Move value from accumulator into register pair H.
func (vm *CPU8080) load_HLA(data []byte) {
	address := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	vm.Logger.Debugf("[77] LD  \t(HL),A ($%04X)", address)
	vm.memory[address] = vm.registers.A
}

// MOV L,A: Load value from accumulator into register L.
func (vm *CPU8080) move_AL(data []byte) {
	vm.Logger.Debugf("[6F] LD  \tL,A")
	vm.registers.A = vm.registers.L
}

// INC H: Increment register pair H.
func (vm *CPU8080) inc_HL(data []byte) {
	vm.Logger.Debugf("[23] INC \tHL")
	hl := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	hl++
	vm.registers.H = byte(hl >> 8)
	vm.registers.L = byte(hl & 0xFF)
}

// INC D: Increment register pair D.
func (vm *CPU8080) inc_DE(data []byte) {
	vm.Logger.Debugf("[13] INC \tDE")
	de := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	de++
	vm.registers.H = byte(de >> 8)
	vm.registers.L = byte(de & 0xFF)
}

// DCR B: Decrement register B.
func (vm *CPU8080) dec_B(data []byte) {
	vm.Logger.Debugf("[05] DEC \tB")
	result := uint16(vm.registers.B) - 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.H = auxCarrySub(vm.registers.B, 1)
	vm.flags.setP(result)

	vm.registers.B--
}

// JNZ addr: Jump if not zero.
func (vm *CPU8080) jump_NZ(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", operand)
	if !vm.flags.Z {
		vm.pc = operand
	} else {
		vm.pc += 2
	}
}

// RET: Return from subroutine.
func (vm *CPU8080) ret(data []byte) {
	address := toUint16(&[]byte{vm.memory[vm.sp], vm.memory[vm.sp+1]})
	vm.Logger.Debugf("[C9] RET \t($%04X)", address)
	vm.pc = address
	vm.sp += 2
}

// MVI HL: Move 8-bit immediate value into address from register pair HL
func (vm *CPU8080) moveI_HL(data []byte) {
	address := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	vm.Logger.Debugf("[36] LD  \t(HL),$%02X", data[0])
	vm.memory[address] = data[0]
	vm.pc++
}

// MOV A,H: Move value from register H into accumulator.
func (vm *CPU8080) move_HA(data []byte) {
	vm.Logger.Debugf("[7E] LD  \tA,H")
	vm.registers.A = vm.registers.H
}

// CPI D8: Compare 8-bit immediate value with accumulator.
func (vm *CPU8080) cmp(data []byte) {
	vm.Logger.Debugf("[FE] CP  \t$%02X", data[0])
	result := uint16(vm.registers.A) - uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = carrySub(vm.registers.A, data[0])
	vm.flags.H = auxCarrySub(vm.registers.A, data[0])
	vm.flags.setP(result)

	vm.registers.A -= data[0]
	vm.pc++
}

// PUSH D: Push register pair D onto stack.
func (vm *CPU8080) push_DE(data []byte) {
	vm.Logger.Debugf("[D5] PUSH\tDE")
	vm.memory[vm.sp-1] = vm.registers.D
	vm.memory[vm.sp-2] = vm.registers.E
	vm.sp -= 2
}

// PUSH H: Push register pair H onto stack.
func (vm *CPU8080) push_HL(data []byte) {
	vm.Logger.Debugf("[E5] PUSH\tHL")
	vm.memory[vm.sp-1] = vm.registers.H
	vm.memory[vm.sp-2] = vm.registers.L
	vm.sp -= 2
}

// DAD H: Add register pair H to register pair H.
func (vm *CPU8080) dad_HL(data []byte) {
	vm.Logger.Debugf("[29] ADD \tHL,HL")
	hl := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	doubledHL := hl << 1

	vm.flags.C = doubledHL > 0xFFFF

	vm.registers.H = byte(doubledHL >> 8)
	vm.registers.L = byte(doubledHL & 0xFF)
}

// DAD D: Add register pair D to register pair H.
func (vm *CPU8080) dad_DE(data []byte) {
	vm.Logger.Debugf("[19] ADD \tHL,DE")
	de := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	doubledDE := de << 1

	vm.flags.C = doubledDE > 0xFFFF

	vm.registers.H = byte(doubledDE >> 8)
	vm.registers.L = byte(doubledDE & 0xFF)
}

// XCHG: Exchange register pairs D and H.
func (vm *CPU8080) xchg(data []byte) {
	vm.Logger.Debugf("[EB] EX  \tDE,HL")
	vm.registers.D, vm.registers.H = vm.registers.H, vm.registers.D
	vm.registers.E, vm.registers.L = vm.registers.L, vm.registers.E
}

// POP H: Pop register pair H from stack.
func (vm *CPU8080) pop_HL(data []byte) {
	vm.Logger.Debugf("[F1] POP \tHL")
	vm.registers.L = vm.memory[vm.sp]
	vm.registers.H = vm.memory[vm.sp+1]
	vm.sp += 2
}
