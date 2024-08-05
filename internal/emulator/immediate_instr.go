package emulator

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) load_SP(data []byte) {
	operand := toUint16(data[1], data[0])
	vm.Logger.Debugf("[31] LD  \tSP,$%04X", operand)
	vm.sp = operand
	vm.PC += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) load_BC(data []byte) {
	vm.Logger.Debugf("[01] LD  \tB,$%04X", toUint16(data[1], data[0]))
	vm.Registers.C = data[0]
	vm.Registers.B = data[1]
	vm.PC += 2
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) load_DE(data []byte) {
	vm.Logger.Debugf("[11] LD  \tDE,$%04X", toUint16(data[1], data[0]))
	vm.Registers.E = data[0]
	vm.Registers.D = data[1]
	vm.PC += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) load_HL(data []byte) {
	vm.Logger.Debugf("[21] LD  \tHL,$%04X", toUint16(data[1], data[0]))
	vm.Registers.L = data[0]
	vm.Registers.H = data[1]
	vm.PC += 2
}

// MVI A, D8: Move 8-bit immediate value into accumulator.
func (vm *CPU8080) moveImm_A(data []byte) {
	vm.Logger.Debugf("[3E] LD  \tA,$%02X", data[0])
	vm.Registers.A = data[0]
	vm.PC++
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveImm_B(data []byte) {
	vm.Logger.Debugf("[06] LD  \tB,$%02X", data[0])
	vm.Registers.B = data[0]
	vm.PC++
}

// MVI C, D8: Move 8-bit immediate value into register C.
func (vm *CPU8080) moveImm_C(data []byte) {
	vm.Logger.Debugf("[0E] LD  \tC,$%02X", data[0])
	vm.Registers.C = data[0]
	vm.PC++
}

// MVI E, D8: Move 8-bit immediate value into register E.
func (vm *CPU8080) moveImm_E(data []byte) {
	vm.Logger.Debugf("[1E] LD  \tE,$%02X", data[0])
	vm.Registers.E = data[0]
	vm.PC++
}

// MVI H, D8: Move 8-bit immediate value into register H.
func (vm *CPU8080) moveImm_H(data []byte) {
	vm.Logger.Debugf("[26] LD  \tH,$%02X", data[0])
	vm.Registers.H = data[0]
	vm.PC++
}

// MVI L, D8: Move 8-bit immediate value into register L.
func (vm *CPU8080) moveImm_L(data []byte) {
	vm.Logger.Debugf("[2E] LD  \tL,$%02X", data[0])
	vm.Registers.L = data[0]
	vm.PC++
}

// MVI D, D8: Move 8-bit immediate value into register L.
func (vm *CPU8080) moveImm_D(data []byte) {
	vm.Logger.Debugf("[16] LD  \tD,$%02X", data[0])
	vm.Registers.D = data[0]
	vm.PC++
}

// MVI M: Move 8-bit immediate value into memory address from register pair HL
func (vm *CPU8080) moveImm_M(data []byte) {
	address := toUint16(vm.Registers.H, vm.Registers.L)
	vm.Logger.Debugf("[36] LD  \t(HL),$%02X", data[0])
	vm.Memory[address] = data[0]
	vm.PC++
}

// ADI: ADD accumulator with 8-bit immediate value.
func (vm *CPU8080) adi(data []byte) {
	vm.Logger.Debugf("[C6] ADD \tA,$%02X", data[0])

	vm.Registers.A = vm.add(data[0])
	vm.PC++
}

// ACI: ADD accumulator with 8-bit immediate value with carry.
func (vm *CPU8080) aci(data []byte) {
	carry := byte(0)
	if vm.flags.C {
		carry = 1
	}
	vm.Registers.A = vm.add(data[0] + carry)

	vm.PC++
}

// SUI: Subtract immediate value from accumulator.
func (vm *CPU8080) sui(data []byte) {
	vm.Logger.Debugf("[D6] SUB \t$%02X", data[0])

	vm.Registers.A = vm.sub(data[0])
	vm.PC++
}

// SBI: Subtract immediate value from accumulator with borrow.
func (vm *CPU8080) sbi(data []byte) {
	vm.Logger.Debugf("[DE] SBI \t$%02X", data[0])
	carry := byte(0)
	if vm.flags.C {
		carry = 1
	}
	subtrahend := uint16(data[0]) + uint16(carry)
	result := uint16(vm.Registers.A) - subtrahend

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.setP(result)
	vm.flags.C = result > 0xFF

	vm.flags.H = (int(vm.Registers.A&0x0F) - int(data[0]&0x0F) - int(carry)) < 0

	vm.Registers.A = byte(result)
	vm.PC++
}

// XRI: Exclusive OR immediate value with accumulator.
func (vm *CPU8080) xri(data []byte) {
	vm.Logger.Debug("[EE] XRI \t$%02X", data[0])
	vm.xra(data[0])
	vm.PC++
}

// ORI: OR A with immediate 8bit value
func (vm *CPU8080) ori(data []byte) {
	vm.Logger.Debugf("[F6] ORI \t$%02X", data[0])
	vm.ora(data[0])
	vm.PC++
}

// ANI D8: AND accumulator with 8-bit immediate value.
func (vm *CPU8080) and(data []byte) {
	vm.Logger.Debugf("[E6] AND \t$%02X", data[0])
	result := uint16(vm.Registers.A) & uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.Registers.A = byte(result)
	vm.PC++
}

// CPI D8: Compare 8-bit immediate value with accumulator.
func (vm *CPU8080) cmp(data []byte) {
	vm.Logger.Debugf("[FE] CP  \t$%02X", data[0])

	vm.compare(data[0])
	vm.PC++
}
