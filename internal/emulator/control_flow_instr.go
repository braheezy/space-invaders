package emulator

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[C3] JMP \t$%04X", operand)
	vm.pc = operand
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

// JZ addr: Jump if zero.
func (vm *CPU8080) jump_Z(data []byte) {
	operand := toUint16(&data)
	vm.Logger.Debugf("[CA] JP  \tZ,$%04X", operand)
	if vm.flags.Z {
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
