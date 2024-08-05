package emulator

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
