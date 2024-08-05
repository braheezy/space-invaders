package emulator

// PCHL: Load program counter from H and L registers.
func (vm *CPU8080) pchl(data []byte) {
	vm.PC = toUint16(vm.Registers.H, vm.Registers.L)
	vm.Logger.Debugf("[E9] PCHL\t($%04X)", vm.PC)
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
