package emulator

// call helper
func (vm *CPU8080) _call(jumpAddress uint16) {
	returnAddress := vm.PC + 2
	vm.push(byte(returnAddress&0xFF), byte(returnAddress>>8))
	vm.PC = jumpAddress
}

// CALL addr: Call subroutine at address
func (vm *CPU8080) call(data []byte) {
	jumpAddress := toUint16(data[1], data[0])
	vm.Logger.Debugf("[CD] CALL\t$%04X", jumpAddress)
	vm._call(jumpAddress)
}

// CNZ addr: Call subroutine at address if zero flag not set
func (vm *CPU8080) call_NZ(data []byte) {
	if !vm.flags.Z {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[C4] CALL\tNZ,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[C4] CALL\tNZ,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CZ addr: Call subroutine at address if zero flag is set
func (vm *CPU8080) call_Z(data []byte) {
	if vm.flags.Z {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[CC] CALL\tZ,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[CC] CALL\tZ,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CC addr: Call subroutine at address if carry flag set
func (vm *CPU8080) call_C(data []byte) {
	if vm.flags.C {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[DC] CALL\tC,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[DC] CALL\tC,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CNC addr: Call subroutine at address if carry flag not set
func (vm *CPU8080) call_NC(data []byte) {
	if !vm.flags.C {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[D4] CALL\tNC,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[D4] CALL\tNC,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CP addr: Call subroutine at address if plus (sign flag is not set)
func (vm *CPU8080) call_P(data []byte) {
	if !vm.flags.S {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[F4] CALL\tP,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[F4] CALL\tP,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CM addr: Call subroutine at address if minus (sign flag set)
func (vm *CPU8080) call_M(data []byte) {
	if vm.flags.S {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[FC] CALL\tM,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[FC] CALL\tM,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CPO addr: Call subroutine at address if parity is odd (not set)
func (vm *CPU8080) call_PO(data []byte) {
	if !vm.flags.P {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[E4] CALL\tPO,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[E4] CALL\tPO,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// CPE addr: Call subroutine at address if parity is even (set)
func (vm *CPU8080) call_PE(data []byte) {
	if vm.flags.P {
		jumpAddress := toUint16(data[1], data[0])
		vm.Logger.Debugf("[EC] CALL\tPE,$%04X", jumpAddress)
		vm._call(jumpAddress)
	} else {
		vm.Logger.Debugf("[EC] CALL\tPE,$%04X (not taken)", vm.PC+2)
		vm.PC += 2
	}
}

// JMP: Jump to address.
func (vm *CPU8080) jump(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[C3] JMP \t$%04X", address)
	vm.PC = address
}

// JNZ addr: Jump if not zero.
func (vm *CPU8080) jump_NZ(data []byte) {
	address := toUint16(data[1], data[0])
	if !vm.flags.Z {
		vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[C2] JP  \tNZ,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JZ addr: Jump if zero.
func (vm *CPU8080) jump_Z(data []byte) {
	address := toUint16(data[1], data[0])
	if vm.flags.Z {
		vm.Logger.Debugf("[CA] JP  \tZ,$%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[CA] JP  \tZ,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JNC addr: Jump if not carry.
func (vm *CPU8080) jump_NC(data []byte) {
	address := toUint16(data[1], data[0])
	if !vm.flags.C {
		vm.Logger.Debugf("[D2] JP  \tNC, $%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[D2] JP  \tNC,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JC addr: Jump if carry.
func (vm *CPU8080) jump_C(data []byte) {
	address := toUint16(data[1], data[0])
	if vm.flags.C {
		vm.Logger.Debugf("[DA] JP  \tC, $%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[DA] JP  \tC,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JM addr: Jump if minus.
func (vm *CPU8080) jump_M(data []byte) {
	address := toUint16(data[1], data[0])
	if vm.flags.S {
		vm.Logger.Debugf("[FA] JP  \tM, $%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[FA] JP  \tM,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JPE addr: Jump if parity is even.
func (vm *CPU8080) jump_PE(data []byte) {
	address := toUint16(data[1], data[0])
	if vm.flags.P {
		vm.Logger.Debugf("[EA] JP  \tPE, $%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[EA] JP  \tPE,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JPO addr: Jump if parity is odd.
func (vm *CPU8080) jump_PO(data []byte) {
	address := toUint16(data[1], data[0])
	if !vm.flags.P {
		vm.Logger.Debugf("[E2] JP  \tPO,$%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[E2] JP  \tPO,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// JP addr: Jump if plus (sign bit is not set).
func (vm *CPU8080) jump_P(data []byte) {
	address := toUint16(data[1], data[0])
	if !vm.flags.S {
		vm.Logger.Debugf("[F2] JP  \tP,$%04X", address)
		vm.PC = address
	} else {
		vm.Logger.Debugf("[F2] JP  \tP,$%04X", vm.PC+2)
		vm.PC += 2
	}
}

// return helper
func (vm *CPU8080) _ret() {
	address := toUint16(vm.Memory[vm.sp+1], vm.Memory[vm.sp])
	vm.PC = address
	vm.sp += 2
}

// RET: Return from subroutine.
func (vm *CPU8080) ret(data []byte) {
	vm._ret()
	vm.Logger.Debugf("[C9] RET \t($%04X)", vm.PC)
}

// RZ: Return from subroutine if Z flag is set.
func (vm *CPU8080) ret_Z(data []byte) {
	if vm.flags.Z {
		vm._ret()
		vm.Logger.Debugf("[C8] RET \tZ($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[C8] RET \tZ (not taken)")
	}
}

// RNZ: Return from subroutine if Z flag is not set.
func (vm *CPU8080) ret_NZ(data []byte) {
	if !vm.flags.Z {
		vm._ret()
		vm.Logger.Debugf("[C0] RET \tNZ($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[C0] RET \tNZ (not taken)")
	}
}

// RC: Return from subroutine if C flag is set.
func (vm *CPU8080) ret_C(data []byte) {
	if vm.flags.C {
		vm._ret()
		vm.Logger.Debugf("[D8] RET \tC($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[D8] RET \tC (not taken)")
	}
}

// RNC: Return from subroutine if C flag is not set.
func (vm *CPU8080) ret_NC(data []byte) {
	if !vm.flags.C {
		vm._ret()
		vm.Logger.Debugf("[D0] RET \tNC($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[D0] RET \tNC (not taken)")
	}
}

// RPE: Return from subroutine if parity even (is set)
func (vm *CPU8080) ret_PE(data []byte) {
	if vm.flags.P {
		vm._ret()
		vm.Logger.Debugf("[E8] RET \tPE($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[E8] RET \tPE (not taken)")
	}
}

// RPO: Return from subroutine if parity odd (is not set)
func (vm *CPU8080) ret_PO(data []byte) {
	if !vm.flags.P {
		vm._ret()
		vm.Logger.Debugf("[E0] RET \tPO($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[E0] RET \tPO (not taken)")
	}
}

// RP: Return from subroutine if plus (sign is not set)
func (vm *CPU8080) ret_P(data []byte) {
	if !vm.flags.S {
		vm._ret()
		vm.Logger.Debugf("[F0] RET \tP($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[F0] RET \tP (not taken)")
	}
}

// RP: Return from subroutine if minus (sign is set)
func (vm *CPU8080) ret_M(data []byte) {
	if vm.flags.S {
		vm._ret()
		vm.Logger.Debugf("[F8] RET \tM($%04X)", vm.PC)
	} else {
		vm.Logger.Debugf("[F8] RET \tM (not taken)")
	}
}

// PCHL: Load program counter from H and L registers.
func (vm *CPU8080) pchl(data []byte) {
	vm.PC = toUint16(vm.Registers.H, vm.Registers.L)
	vm.Logger.Debugf("[E9] PCHL\t($%04X)", vm.PC)
}
