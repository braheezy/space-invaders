package emulator

func (vm *CPU8080) nextOpCode() {
	currentCode := (*vm.programData)[vm.pc : vm.pc+3]

	op := currentCode[0]
	vm.pc += 1

	switch op {
	case 0x00:
		// NOP: No operation.
		vm.Logger.Debugf("[%02X] NOP", op)
	case 0x01:
		// LXI: Load 16-bit immediate value into register pair B,C
		vm.registers.C = currentCode[1]
		vm.registers.B = currentCode[2]
		vm.pc += 2
	case 0xC3:
		// JMP: Jump to address.
		operand := getOperand(&currentCode)
		vm.Logger.Debugf("[%02X] JMP to %04X", op, operand)
		vm.pc = operand
	default:
		// Unsupported opcode
		vm.Logger.Fatalf("Unsupported opcode: %02X", op)
	}
}

func getOperand(code *[]byte) uint16 {
	return uint16((*code)[2])<<8 | uint16((*code)[1])
}
