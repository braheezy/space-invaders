package emulator

// xra performs Exclusive OR register with accumulator
func (vm *CPU8080) xra(reg byte) {
	result := uint16(vm.Registers.A) ^ uint16(reg)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)
	vm.flags.H = false

	vm.Registers.A = byte(result)
}

// XRA A: Exclusive-OR accumulator with accumulator.
func (vm *CPU8080) xra_A(data []byte) {
	vm.Logger.Debugf("[AF] XOR \tA")
	vm.xra(vm.Registers.A)
}

// ora performs OR with accumulator
func (vm *CPU8080) ora(reg byte) {
	result := uint16(vm.Registers.A) | uint16(reg)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.Registers.A = byte(result)
}

// ORA B: OR A with register B
func (vm *CPU8080) ora_B(data []byte) {
	vm.Logger.Debugf("[B0] OR  \tB")
	vm.ora(vm.Registers.B)
}

// ORA H: OR A with register H
func (vm *CPU8080) ora_H(data []byte) {
	vm.Logger.Debugf("[B4] OR  \tH")
	vm.ora(vm.Registers.H)
}

// ORA M: OR A with memory location pointed to by register pair HL
func (vm *CPU8080) ora_M(data []byte) {
	vm.Logger.Debugf("[B6] OR  \t(HL)")
	address := toUint16(vm.Registers.H, vm.Registers.L)
	vm.ora(vm.Memory[address])
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

// RRC: Rotate accumulator right.
// The carry bit is set equal to the low-order
// bit of the accumulator. The contents of the accumulator are
// rotated one bit position to the right, with the low-order bit
// being transferred to the high-order bit position of the
// accumulator.
func (vm *CPU8080) rrc(data []byte) {
	vm.Logger.Debugf("[0F] RRC \tA")
	// Isolate least significant bit to check for Carry
	vm.flags.C = vm.Registers.A&0x01 == 1
	// Rotate accumulator right
	vm.Registers.A = (vm.Registers.A >> 1) | (vm.Registers.A << (8 - 1))
}

// RLC: Rotate accumulator left. The Carry bit is set equal to the high-order
// bit of the accumulator. The contents of the accumulator are rotated one bit
// position to the left, with the high-order bit being transferred to the
// low-order bit position of the accumulator
func (vm *CPU8080) rlc(data []byte) {
	vm.Logger.Debugf("[07] RLC \tA")
	// Isolate most significant bit to check for Carry
	vm.flags.C = (vm.Registers.A & 0x80) == 0x80
	// Rotate accumulator left
	vm.Registers.A = (vm.Registers.A << 1) | (vm.Registers.A >> (8 - 1))
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
	vm.flags.C = vm.Registers.A&0x01 != 0
	// Rotate accumulator right through carry
	vm.Registers.A = (vm.Registers.A >> 1) | (carryRotate << (8 - 1))
}

// compare helper
func (vm *CPU8080) compare(data byte) {
	result := uint16(vm.Registers.A) - uint16(data)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.setP(result)
	vm.flags.C = carrySub(vm.Registers.A, data)
	vm.flags.H = auxCarrySub(vm.Registers.A, data)
}

// CPI D8: Compare 8-bit immediate value with accumulator.
func (vm *CPU8080) cmp(data []byte) {
	vm.Logger.Debugf("[FE] CP  \t$%02X", data[0])

	vm.compare(data[0])
	vm.PC++
}

// CMP B: Compare A with register B
func (vm *CPU8080) cmp_B(data []byte) {
	vm.Logger.Debugf("[B8] CP  \tB")
	vm.compare(vm.Registers.B)
}

// CMP M: Compare A with memory address pointed to by register pair HL
func (vm *CPU8080) cmp_M(data []byte) {
	vm.Logger.Debugf("[BE] CP  \t(HL)")
	vm.compare(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// ana performs AND with data and accumulator, storing in accumulator.
func (vm *CPU8080) ana(data byte) {
	result := uint16(vm.Registers.A) & uint16(data)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)

	vm.Registers.A = byte(result)
}

// ANA A: AND accumulator with accumulator.
func (vm *CPU8080) ana_A(data []byte) {
	vm.Logger.Debugf("[A7] AND \tA")
	vm.ana(vm.Registers.A)
}

// ANA B: AND register B with accumulator.
func (vm *CPU8080) ana_B(data []byte) {
	vm.Logger.Debugf("[A0] AND \tB")
	vm.ana(vm.Registers.B)
}

// XTHL: Exchange top of stack with address referenced by register pair HL.
func (vm *CPU8080) xthl(data []byte) {
	vm.Logger.Debugf("[E3] EX  \t(SP),HL")
	vm.Memory[vm.sp], vm.Memory[vm.sp+1] = vm.Registers.L, vm.Registers.H
}

// XCHG: Exchange register pairs D and H.
func (vm *CPU8080) xchg(data []byte) {
	vm.Logger.Debugf("[EB] EX  \tDE,HL")
	vm.Registers.D, vm.Registers.H = vm.Registers.H, vm.Registers.D
	vm.Registers.E, vm.Registers.L = vm.Registers.L, vm.Registers.E
}

// DAA: Decimal Adjust Accumulator
// The eight bit hex number in the accumulator is adjusted to form two
// four bit binary decimal digits.
func (vm *CPU8080) daa(data []byte) {
	vm.Logger.Debugf("[27] DAA")
	// Step 1: Adjust lower nibble
	lower := vm.Registers.A & 0x0F
	if lower > 9 || vm.flags.H {
		vm.Registers.A += 6
		vm.flags.H = lower < 6
	} else {
		vm.flags.H = false
	}

	// Step 2: Adjust upper nibble
	if vm.Registers.A > 0x9F || vm.flags.C {
		vm.Registers.A += 0x60
		vm.flags.C = true
	}

	// Set Zero flag
	vm.flags.Z = (vm.Registers.A == 0)

	// Set Sign flag
	vm.flags.S = (vm.Registers.A & 0x80) != 0

	// Set Parity flag
	vm.flags.P = (vm.Registers.A & 0x01) == 0

	// Corrected Parity flag calculation
	count := 0
	for i := 0; i < 8; i++ {
		if (vm.Registers.A & (1 << i)) != 0 {
			count++
		}
	}
	vm.flags.P = (count % 2) == 0
}
