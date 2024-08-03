package emulator

// xra performs Exclusive OR register with accumulator
func (vm *CPU8080) xra(data byte) {
	result := uint16(vm.Registers.A) ^ uint16(data)

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = false
	vm.flags.setP(result)
	vm.flags.H = false

	vm.Registers.A = byte(result)
}

// XRI: Exclusive OR immediate value with accumulator.
func (vm *CPU8080) xri(data []byte) {
	vm.Logger.Debug("[EE] XRI \t$%02X", data[0])
	vm.xra(data[0])
	vm.PC++
}

// XRA B: Exclusive-OR register B with accumulator.
func (vm *CPU8080) xra_B(data []byte) {
	vm.Logger.Debug("[A8] XOR \tB")
	vm.xra(vm.Registers.B)
}

// XRA C: Exclusive-OR register C with accumulator.
func (vm *CPU8080) xra_C(data []byte) {
	vm.Logger.Debug("[A9] XOR \tC")
	vm.xra(vm.Registers.C)
}

// XRA D: Exclusive-OR register D with accumulator.
func (vm *CPU8080) xra_D(data []byte) {
	vm.Logger.Debug("[AA] XOR \tD")
	vm.xra(vm.Registers.D)
}

// XRA E: Exclusive-OR register E with accumulator.
func (vm *CPU8080) xra_E(data []byte) {
	vm.Logger.Debug("[AB] XOR \tE")
	vm.xra(vm.Registers.E)
}

// XRA H: Exclusive-OR register H with accumulator.
func (vm *CPU8080) xra_H(data []byte) {
	vm.Logger.Debug("[AC] XOR \tH")
	vm.xra(vm.Registers.H)
}

// XRA L: Exclusive-OR register L with accumulator.
func (vm *CPU8080) xra_L(data []byte) {
	vm.Logger.Debug("[AD] XOR \tL")
	vm.xra(vm.Registers.L)
}

// XRA M: Exclusive-OR memory address pointed to by register pair HL with accumulator.
func (vm *CPU8080) xra_M(data []byte) {
	vm.Logger.Debug("[AE] XOR \tM")
	vm.xra(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// XRA A: Exclusive-OR accumulator with accumulator.
func (vm *CPU8080) xra_A(data []byte) {
	vm.Logger.Debug("[AF] XOR \tA")
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
	vm.Logger.Debug("[B0] OR  \tB")
	vm.ora(vm.Registers.B)
}

// ORA C: OR A with register C
func (vm *CPU8080) ora_C(data []byte) {
	vm.Logger.Debug("[B1] OR  \tC")
	vm.ora(vm.Registers.C)
}

// ORA D: OR A with register D
func (vm *CPU8080) ora_D(data []byte) {
	vm.Logger.Debug("[B2] OR  \tD")
	vm.ora(vm.Registers.D)
}

// ORA E: OR A with register E
func (vm *CPU8080) ora_E(data []byte) {
	vm.Logger.Debug("[B3] OR  \tE")
	vm.ora(vm.Registers.E)
}

// ORA H: OR A with register H
func (vm *CPU8080) ora_H(data []byte) {
	vm.Logger.Debug("[B4] OR  \tH")
	vm.ora(vm.Registers.H)
}

// ORA L: OR A with register L
func (vm *CPU8080) ora_L(data []byte) {
	vm.Logger.Debug("[B5] OR  \tL")
	vm.ora(vm.Registers.L)
}

// ORA A: OR A with register A
func (vm *CPU8080) ora_A(data []byte) {
	vm.Logger.Debug("[B7] OR  \tA")
	vm.ora(vm.Registers.A)
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

// RAL: Rotate accumulator left through carry.
// The contents of the accumulator are rotated one bit position to the left.
// The high-order bit of the accumulator replaces the Carry bit, while the
// Carry bit replaces the high-order bit of the accumulator.
func (vm *CPU8080) ral(data []byte) {
	vm.Logger.Debugf("[17] RAL \tA")
	var carry uint8
	if vm.flags.C {
		carry = 1
	}
	// Isolate most significant bit to check for Carry
	vm.flags.C = (vm.Registers.A & 0x80) == 0x80
	// Rotate accumulator left through carry
	vm.Registers.A = (vm.Registers.A << 1) | carry
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

// CMP C: Compare A with register C
func (vm *CPU8080) cmp_C(data []byte) {
	vm.Logger.Debugf("[B9] CP  \tC")
	vm.compare(vm.Registers.C)
}

// CMP D: Compare A with register D
func (vm *CPU8080) cmp_D(data []byte) {
	vm.Logger.Debugf("[BA] CP  \tD")
	vm.compare(vm.Registers.D)
}

// CMP E: Compare A with register E
func (vm *CPU8080) cmp_E(data []byte) {
	vm.Logger.Debugf("[BB] CP  \tE")
	vm.compare(vm.Registers.E)
}

// CMP H: Compare A with register H
func (vm *CPU8080) cmp_H(data []byte) {
	vm.Logger.Debugf("[BC] CP  \tH")
	vm.compare(vm.Registers.H)
}

// CMP L: Compare A with register L
func (vm *CPU8080) cmp_L(data []byte) {
	vm.Logger.Debugf("[BD] CP  \tL")
	vm.compare(vm.Registers.L)
}

// CMP A: Compare A with register A
func (vm *CPU8080) cmp_A(data []byte) {
	vm.Logger.Debugf("[BF] CP  \tA")
	vm.compare(vm.Registers.A)
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

// ANA C: AND register C with accumulator.
func (vm *CPU8080) ana_C(data []byte) {
	vm.Logger.Debugf("[A1] AND \tC")
	vm.ana(vm.Registers.C)
}

// ANA D: AND register D with accumulator.
func (vm *CPU8080) ana_D(data []byte) {
	vm.Logger.Debugf("[A2] AND \tD")
	vm.ana(vm.Registers.D)
}

// ANA E: AND register E with accumulator.
func (vm *CPU8080) ana_E(data []byte) {
	vm.Logger.Debugf("[A3] AND \tE")
	vm.ana(vm.Registers.E)
}

// ANA H: AND register H with accumulator.
func (vm *CPU8080) ana_H(data []byte) {
	vm.Logger.Debugf("[A4] AND \tH")
	vm.ana(vm.Registers.H)
}

// ANA L: AND register L with accumulator.
func (vm *CPU8080) ana_L(data []byte) {
	vm.Logger.Debugf("[A5] AND \tL")
	vm.ana(vm.Registers.L)
}

// ANA M: AND memory address pointed to by register pair HL with accumulator.
func (vm *CPU8080) ana_M(data []byte) {
	vm.Logger.Debug("[A6] AND \tL")
	vm.ana(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
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

// CMA: Complement accumulator.
func (vm *CPU8080) cma(data []byte) {
	vm.Logger.Debugf("[2F] CMA")
	vm.Registers.A = ^vm.Registers.A
}
