package emulator

func (vm *CPU8080) nextOpCode() {
	currentCode := vm.memory[vm.pc : vm.pc+3]

	op := currentCode[0]
	vm.pc += 1

	if opcodeFunc, exists := vm.opcodeTable[op]; exists {
		opcodeFunc(currentCode[1:])
	} else {
		vm.Logger.Fatalf("Unsupported opcode: %02X", op)
	}

}

func toUint16(code *[]byte) uint16 {
	return uint16((*code)[1])<<8 | uint16((*code)[0])
}

// NOP: No operation.
func (vm *CPU8080) nop(data []byte) {
	vm.Logger.Debugf("[00] NOP")
}

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[C3] JMP to $%04X", operand)
	vm.pc = operand
}

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) loadSP(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[31] LXI SP, $%04X", operand)
	vm.sp = operand
	vm.pc += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) loadBC(data []byte) {
	vm.Logger.Debugf("[01] LXI B, $%04X", toUint16(&data))
	vm.registers.C = data[0]
	vm.registers.B = data[1]
	vm.pc += 2
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveB(data []byte) {
	vm.Logger.Debugf("[06] MVI B, $%02X", data[0])
	vm.registers.B = data[0]
	vm.pc += 1
}

// CALL addr: Call subroutine at address
func (vm *CPU8080) call(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[CD] CALL $%04X", operand)
	vm.pc = operand
	vm.memory[vm.sp-1] = data[1]
	vm.memory[vm.sp-2] = data[0]
	vm.sp -= 2
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) loadDE(data []byte) {
	vm.Logger.Debugf("[11] LXI D, $%04X", toUint16(&data))
	vm.registers.E = data[0]
	vm.registers.D = data[1]
	vm.pc += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) loadHL(data []byte) {
	vm.Logger.Debugf("[21] LXI H, $%04X", toUint16(&data))
	vm.registers.L = data[0]
	vm.registers.H = data[1]
	vm.pc += 2
}

// LDAX D: Load value from address in register pair D into accumulator.
func (vm *CPU8080) loadAXD(data []byte) {
	address := toUint16(&[]byte{vm.registers.D, vm.registers.E})
	vm.Logger.Debugf("[1A] LDAX D")
	vm.registers.A = vm.memory[address]
}

// MOV M, A: Move value from accumulator into register pair H.
func (vm *CPU8080) moveMA(data []byte) {
	address := toUint16(&[]byte{vm.registers.H, vm.registers.L})
	vm.Logger.Debugf("[22] MOV M, A")
	vm.memory[address] = vm.registers.A
}
