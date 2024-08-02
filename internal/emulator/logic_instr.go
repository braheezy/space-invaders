package emulator

// xra performs Exclusive OR register with accumulator
func (vm *CPU8080) xra(reg byte) {
	result := uint16(vm.registers.A) ^ uint16(reg)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.registers.A = byte(result)
}

// XRA A: Exclusive-OR accumulator with accumulator.
func (vm *CPU8080) xra_A(data []byte) {
	vm.Logger.Debugf("[AF] XOR \tA")
	vm.xra(vm.registers.A)
}

// ora performs OR with accumulator
func (vm *CPU8080) ora(reg byte) {
	result := uint16(vm.registers.A) | uint16(reg)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.registers.A = byte(result)
}

// ORA B: OR A with register B
func (vm *CPU8080) ora_B(data []byte) {
	vm.Logger.Debugf("[B0] OR  \tB")
	vm.ora(vm.registers.B)
}

// ORA H: OR A with register H
func (vm *CPU8080) ora_H(data []byte) {
	vm.Logger.Debugf("[B4] OR  \tH")
	vm.ora(vm.registers.H)
}

// ORA M: OR A with memory location pointed to by register pair HL
func (vm *CPU8080) ora_M(data []byte) {
	vm.Logger.Debugf("[B6] OR  \t(HL)")
	address := toUint16(vm.registers.H, vm.registers.L)
	vm.ora(vm.memory[address])
}

// ORI: OR A with immediate 8bit value
func (vm *CPU8080) ori(data []byte) {
	vm.Logger.Debugf("[F6] ORI \t$%02X", data[0])
	vm.ora(data[0])
	vm.pc++
}

// ANI D8: AND accumulator with 8-bit immediate value.
func (vm *CPU8080) and(data []byte) {
	vm.Logger.Debugf("[E6] AND \t$%02X", data[0])
	result := uint16(vm.registers.A) & uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.registers.A = byte(result)
	vm.pc++
}

// RRC: Rotate accumulator right.
// The carry bit is set equal to the low-order
// bit of the accumulator. The contents of the accumulator are
// rotated one bit position to the right, with the low-order bit
// being transferred to the high-order bit position of the
// accumulator.
func (vm *CPU8080) rrc(data []byte) {
	vm.Logger.Debugf("[0F] RRC \tA")
	// Isolate least significant bit to check for Carry
	vm.flags.C = vm.registers.A&0x01 == 1
	// Rotate accumulator right
	vm.registers.A = (vm.registers.A >> 1) | (vm.registers.A << (8 - 1))
}

// RLC: Rotate accumulator left. The Carry bit is set equal to the high-order
// bit of the accumulator. The contents of the accumulator are rotated one bit
// position to the left, with the high-order bit being transferred to the
// low-order bit position of the accumulator
func (vm *CPU8080) rlc(data []byte) {
	vm.Logger.Debugf("[07] RLC \tA")
	// Isolate most significant bit to check for Carry
	vm.flags.C = (vm.registers.A & 0x80) == 0x80
	// Rotate accumulator left
	vm.registers.A = (vm.registers.A << 1) | (vm.registers.A >> (8 - 1))
}

// RAR: Rotate accumulator right through carry.
// The contents of the accumulator are rotated one bit position to the right.
// The low order bit of the accumulator replaces the carry bit, while the carry bit replaces
// the high order bit of the accumulator.
func (vm *CPU8080) rar(data []byte) {
	vm.Logger.Debugf("[1F] RAR \tA")
	var carryRotate uint8
	if vm.flags.C {
		carryRotate = 1
	}
	// Isolate least significant bit to check for Carry
	vm.flags.C = vm.registers.A&0x01 != 0
	// Rotate accumulator right through carry
	vm.registers.A = (vm.registers.A >> 1) | (carryRotate << (8 - 1))
}

// RAL: Rotate accumulator left through carry.
// The contents of the accumulator are rotated one bit position to the left.
// The high order bit of the accumulator replaces the carry bit, while the carry bit replaces
// the low order bit of the accumulator.

// CPI D8: Compare 8-bit immediate value with accumulator.
func (vm *CPU8080) cmp(data []byte) {
	vm.Logger.Debugf("[FE] CP  \t$%02X", data[0])
	result := uint16(vm.registers.A) - uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = carrySub(vm.registers.A, data[0])
	vm.flags.H = auxCarrySub(vm.registers.A, data[0])
	vm.flags.setP(result)

	vm.pc++
}

// ana performs AND with data and accumulator, storing in accumulator.
func (vm *CPU8080) ana(data byte) {
	result := uint16(vm.registers.A) & uint16(data)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.registers.A = byte(result)
}

// ANA A: AND accumulator with accumulator.
func (vm *CPU8080) ana_A(data []byte) {
	vm.Logger.Debugf("[A7] AND \tA")
	vm.ana(vm.registers.A)
}

// XTHL: Exchange top of stack with address referenced by register pair HL.
func (vm *CPU8080) xthl(data []byte) {
	vm.Logger.Debugf("[E3] EX  \t(SP),HL")
	vm.memory[vm.sp], vm.memory[vm.sp+1] = vm.registers.L, vm.registers.H
}

// XCHG: Exchange register pairs D and H.
func (vm *CPU8080) xchg(data []byte) {
	vm.Logger.Debugf("[EB] EX  \tDE,HL")
	vm.registers.D, vm.registers.H = vm.registers.H, vm.registers.D
	vm.registers.E, vm.registers.L = vm.registers.L, vm.registers.E
}

// DAA: Decimal Adjust Accumulator
// The eight bit hex number in the accumulator is adjusted to form two
// four bit binary decimal digits.
func (vm *CPU8080) daa(data []byte) {
	vm.Logger.Debugf("[27] DAA")
	// Step 1: Adjust lower nibble
	lower := vm.registers.A & 0x0F
	if lower > 9 || vm.flags.H {
		vm.registers.A += 6
		vm.flags.H = lower < 6
	} else {
		vm.flags.H = false
	}

	// Step 2: Adjust upper nibble
	if vm.registers.A > 0x9F || vm.flags.C {
		vm.registers.A += 0x60
		vm.flags.C = true
	}

	// Set Zero flag
	vm.flags.Z = (vm.registers.A == 0)

	// Set Sign flag
	vm.flags.S = (vm.registers.A & 0x80) != 0

	// Set Parity flag
	vm.flags.P = (vm.registers.A & 0x01) == 0

	// Corrected Parity flag calculation
	count := 0
	for i := 0; i < 8; i++ {
		if (vm.registers.A & (1 << i)) != 0 {
			count++
		}
	}
	vm.flags.P = (count % 2) == 0
}
