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
