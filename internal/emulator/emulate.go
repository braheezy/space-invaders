package emulator

func (vm *CPU8080) nextOpCode() {
	currentCode := (*vm.programData)[vm.pc : vm.pc+3]

	op := currentCode[0]
	vm.pc += 1

	if opcodeFunc, exists := vm.opcodeTable[op]; exists {
		opcodeFunc(currentCode[1:])
	} else {
		vm.Logger.Fatalf("Unsupported opcode: %02X", op)
	}

}

func toOperand(code *[]byte) uint16 {
	return uint16((*code)[1])<<8 | uint16((*code)[0])
}

// NOP: No operation.
func (vm *CPU8080) nop(data []byte) {
	vm.Logger.Debugf("[00] NOP")
}

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	operand := toOperand(&data)
	vm.Logger.Debugf("[C3] JMP to $%04X", operand)
	vm.pc = operand
}

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) loadSP(data []byte) {
	operand := toOperand(&data)
	vm.Logger.Debugf("[31] LXI SP, $%04X", operand)
	vm.sp = operand
	vm.pc += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) loadBC(data []byte) {
	vm.Logger.Debugf("[01] LXI B, $%04X", toOperand(&data))
	vm.registers.C = data[1]
	vm.registers.B = data[2]
	vm.pc += 2
}
