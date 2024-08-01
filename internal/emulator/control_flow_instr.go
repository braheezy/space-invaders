package emulator

// CALL addr: Call subroutine at address
func (vm *CPU8080) call(data []byte) {
	jumpAddress := toUint16(data[1], data[0])
	returnAddress := vm.pc + 2
	vm.Logger.Debugf("[CD] CALL\t$%04X", jumpAddress)
	vm.push(byte(returnAddress&0xFF), byte(returnAddress>>8))
	vm.pc = jumpAddress
}

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	operand := toUint16(data[1], data[0])
	vm.Logger.Debugf("[C3] JMP \t$%04X", operand)
	vm.pc = operand
}

// JNZ addr: Jump if not zero.
func (vm *CPU8080) jump_NZ(data []byte) {
	operand := toUint16(data[1], data[0])
	if !vm.flags.Z {
		vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", operand)
		vm.pc = operand
	} else {
		vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", vm.pc+2)
		vm.pc += 2
	}
}

// JZ addr: Jump if zero.
func (vm *CPU8080) jump_Z(data []byte) {
	operand := toUint16(data[1], data[0])
	if vm.flags.Z {
		vm.Logger.Debugf("[CA] JP  \tZ,$%04X", operand)
		vm.pc = operand
	} else {
		vm.Logger.Debugf("[CA] JP  \tZ,$%04X", vm.pc+2)
		vm.pc += 2
	}
}

// JNC addr: Jump if not carry.
func (vm *CPU8080) jump_NC(data []byte) {
	operand := toUint16(data[1], data[0])
	if !vm.flags.C {
		vm.Logger.Debugf("[D2] JP  \tNC, $%04X", operand)
		vm.pc = operand
	} else {
		vm.Logger.Debugf("[D2] JP  \tNC,$%04X", vm.pc+2)
		vm.pc += 2
	}
}

// JC addr: Jump if carry.
func (vm *CPU8080) jump_C(data []byte) {
	operand := toUint16(data[1], data[0])
	if vm.flags.C {
		vm.Logger.Debugf("[DA] JP  \tC, $%04X", operand)
		vm.pc = operand
	} else {
		vm.Logger.Debugf("[DA] JP  \tC,$%04X", vm.pc+2)
		vm.pc += 2
	}
}

// return helper
func (vm *CPU8080) _ret() {
	address := toUint16(vm.memory[vm.sp+1], vm.memory[vm.sp])
	vm.pc = address
	vm.sp += 2
}

// RET: Return from subroutine.
func (vm *CPU8080) ret(data []byte) {
	vm._ret()
	vm.Logger.Debugf("[C9] RET \t($%04X)", vm.pc)
}

// RZ: Return from subroutine if Z flag is set.
func (vm *CPU8080) ret_Z(data []byte) {
	if vm.flags.Z {
		vm._ret()
		vm.Logger.Debugf("[C8] RET \tZ($%04X)", vm.pc)
	} else {
		vm.Logger.Debugf("[C8] RET \tZ (not taken)")
	}
}

// RNZ: Return from subroutine if Z flag is not set.
func (vm *CPU8080) ret_NZ(data []byte) {
	if !vm.flags.Z {
		vm._ret()
		vm.Logger.Debugf("[C0] RET \tNZ($%04X)", vm.pc)
	} else {
		vm.Logger.Debugf("[C0] RET \tNZ (not taken)")
	}
}

// RC: Return from subroutine if C flag is set.
func (vm *CPU8080) ret_C(data []byte) {
	if vm.flags.C {
		vm._ret()
		vm.Logger.Debugf("[D8] RET \tC($%04X)", vm.pc)
	} else {
		vm.Logger.Debugf("[D8] RET \tC (not taken)")
	}
}

// RNC: Return from subroutine if C flag is not set.
func (vm *CPU8080) ret_NC(data []byte) {
	if !vm.flags.C {
		vm._ret()
		vm.Logger.Debugf("[D0] RET \tNC($%04X)", vm.pc)
	} else {
		vm.Logger.Debugf("[D0] RET \tNC (not taken)")
	}
}
