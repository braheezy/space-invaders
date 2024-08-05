package emulator

// add helper
func (vm *CPU8080) add(data byte) byte {
	result := vm.Registers.A + data

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carryAdd(vm.Registers.A, data)
	vm.flags.H = auxCarryAdd(vm.Registers.A, data)
	vm.flags.setP(uint16(result))

	return byte(result)
}

// ADD A: ADD accumulator with register A.
func (vm *CPU8080) add_A(data []byte) {
	vm.Logger.Debug("[87] ADD \tA,A")

	vm.Registers.A = vm.add(vm.Registers.A)
}

// ADD B: ADD accumulator with register B.
func (vm *CPU8080) add_B(data []byte) {
	vm.Logger.Debug("[80] ADD \tA,B")

	vm.Registers.A = vm.add(vm.Registers.B)
}

// ADD C: ADD accumulator with register C.
func (vm *CPU8080) add_C(data []byte) {
	vm.Logger.Debug("[81] ADD \tA,C")

	vm.Registers.A = vm.add(vm.Registers.C)
}

// ADD D: ADD accumulator with register D.
func (vm *CPU8080) add_D(data []byte) {
	vm.Logger.Debug("[82] ADD \tA,D")

	vm.Registers.A = vm.add(vm.Registers.D)
}

// ADD E: ADD accumulator with register E.
func (vm *CPU8080) add_E(data []byte) {
	vm.Logger.Debug("[83] ADD \tA,E")

	vm.Registers.A = vm.add(vm.Registers.E)
}

// ADD H: ADD accumulator with register H.
func (vm *CPU8080) add_H(data []byte) {
	vm.Logger.Debug("[84] ADD \tA,H")

	vm.Registers.A = vm.add(vm.Registers.H)
}

// ADD L: ADD accumulator with register L.
func (vm *CPU8080) add_L(data []byte) {
	vm.Logger.Debug("[85] ADD \tA,L")

	vm.Registers.A = vm.add(vm.Registers.L)
}

// ADD M: ADD accumulator with memory address pointed to by register pair HL
func (vm *CPU8080) add_M(data []byte) {
	vm.Logger.Debug("[86] ADD \tA,(HL)")

	vm.Registers.A = vm.add(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// add with carry helper
func (vm *CPU8080) adc(data byte) byte {
	carry := byte(0)
	if vm.flags.C {
		carry = 1
	}
	result := vm.Registers.A + data + carry

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carryAdd(vm.Registers.A, data+byte(carry))
	vm.flags.H = ((vm.Registers.A & 0xF) + (data & 0xF) + byte(carry)) > 0xF
	vm.flags.setP(uint16(result))

	return byte(result)
}

// ADC A: Add accumulator with register A and carry.
func (vm *CPU8080) adc_A(data []byte) {
	vm.Logger.Debug("[8F] ADC \tA")

	vm.Registers.A = vm.adc(vm.Registers.A)
}

// ADC B: Add accumulator with register B and carry.
func (vm *CPU8080) adc_B(data []byte) {
	vm.Logger.Debug("[88] ADC \tB")

	vm.Registers.A = vm.adc(vm.Registers.B)
}

// ADC C: Add accumulator with register C and carry.
func (vm *CPU8080) adc_C(data []byte) {
	vm.Logger.Debug("[89] ADC \tC")

	vm.Registers.A = vm.adc(vm.Registers.C)
}

// ADC D: Add accumulator with register D and carry.
func (vm *CPU8080) adc_D(data []byte) {
	vm.Logger.Debug("[8A] ADC \tD")

	vm.Registers.A = vm.adc(vm.Registers.D)
}

// ADC E: Add accumulator with register E and carry.
func (vm *CPU8080) adc_E(data []byte) {
	vm.Logger.Debug("[8B] ADC \tE")

	vm.Registers.A = vm.adc(vm.Registers.E)
}

// ADC H: Add accumulator with register H and carry.
func (vm *CPU8080) adc_H(data []byte) {
	vm.Logger.Debug("[8C] ADC \tH")

	vm.Registers.A = vm.adc(vm.Registers.H)
}

// ADC L: Add accumulator with register L and carry.
func (vm *CPU8080) adc_L(data []byte) {
	vm.Logger.Debug("[8D] ADC \tL")

	vm.Registers.A = vm.adc(vm.Registers.L)
}

// ADC M: Subtract memory address pointed to by register pair HL from accumulator.
func (vm *CPU8080) adc_M(data []byte) {
	vm.Logger.Debug("[8E] ADC \tM")

	vm.Registers.A = vm.adc(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// subtract helper
func (vm *CPU8080) sub(data byte) byte {
	result := uint16(vm.Registers.A) - uint16(data)

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carrySub(vm.Registers.A, data)
	vm.flags.H = auxCarrySub(vm.Registers.A, data)
	vm.flags.setP(uint16(result))

	return byte(result)
}

// SUB A: Subtract accumulator from accumulator.
func (vm *CPU8080) sub_A(data []byte) {
	vm.Logger.Debug("[97] SUB \tA")

	vm.Registers.A = vm.sub(vm.Registers.A)
}

// SUB B: Subtract register B from accumulator.
func (vm *CPU8080) sub_B(data []byte) {
	vm.Logger.Debug("[90] SUB \tB")

	vm.Registers.A = vm.sub(vm.Registers.B)
}

// SUB C: Subtract register C from accumulator.
func (vm *CPU8080) sub_C(data []byte) {
	vm.Logger.Debug("[91] SUB \tC")

	vm.Registers.A = vm.sub(vm.Registers.C)
}

// SUB D: Subtract register D from accumulator.
func (vm *CPU8080) sub_D(data []byte) {
	vm.Logger.Debug("[92] SUB \tD")

	vm.Registers.A = vm.sub(vm.Registers.D)
}

// SUB E: Subtract register E from accumulator.
func (vm *CPU8080) sub_E(data []byte) {
	vm.Logger.Debug("[93] SUB \tE")

	vm.Registers.A = vm.sub(vm.Registers.E)
}

// SUB H: Subtract register H from accumulator.
func (vm *CPU8080) sub_H(data []byte) {
	vm.Logger.Debug("[94] SUB \tH")

	vm.Registers.A = vm.sub(vm.Registers.H)
}

// SUB L: Subtract register L from accumulator.
func (vm *CPU8080) sub_L(data []byte) {
	vm.Logger.Debug("[95] SUB \tL")

	vm.Registers.A = vm.sub(vm.Registers.L)
}

// SUB M: Subtract memory address pointed to by register pair HL from accumulator.
func (vm *CPU8080) sub_M(data []byte) {
	vm.Logger.Debug("[96] SUB \tL")

	vm.Registers.A = vm.sub(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// subtract with borrow helper
func (vm *CPU8080) sbb(data byte) byte {
	carry := byte(0)
	if vm.flags.C {
		carry = 1
	}
	subtrahend := data + carry
	result := vm.Registers.A - subtrahend

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carrySub(vm.Registers.A, subtrahend)
	vm.flags.H = ((vm.Registers.A & 0x0F) - (subtrahend & 0x0F)) != 0
	vm.flags.setP(uint16(result))

	return byte(result)
}

// SBB A: Subtract register A from accumulator with borrow.
func (vm *CPU8080) sbb_A(data []byte) {
	vm.Logger.Debug("[9F] SBB \tA")

	vm.Registers.A = vm.sbb(vm.Registers.A)
}

// SBB B: Subtract register B from accumulator with borrow.
func (vm *CPU8080) sbb_B(data []byte) {
	vm.Logger.Debug("[98] SBB \tB")

	vm.Registers.A = vm.sbb(vm.Registers.B)
}

// SBB C: Subtract register C from accumulator with borrow.
func (vm *CPU8080) sbb_C(data []byte) {
	vm.Logger.Debug("[99] SBB \tC")

	vm.Registers.A = vm.sbb(vm.Registers.C)
}

// SBB D: Subtract register D from accumulator with borrow.
func (vm *CPU8080) sbb_D(data []byte) {
	vm.Logger.Debug("[9A] SBB \tD")

	vm.Registers.A = vm.sbb(vm.Registers.D)
}

// SBB E: Subtract register E from accumulator with borrow.
func (vm *CPU8080) sbb_E(data []byte) {
	vm.Logger.Debug("[9B] SBB \tE")

	vm.Registers.A = vm.sbb(vm.Registers.E)
}

// SBB H: Subtract register H from accumulator with borrow.
func (vm *CPU8080) sbb_H(data []byte) {
	vm.Logger.Debug("[9C] SBB \tH")

	vm.Registers.A = vm.sbb(vm.Registers.H)
}

// SBB L: Subtract register L from accumulator with borrow.
func (vm *CPU8080) sbb_L(data []byte) {
	vm.Logger.Debug("[9D] SBB \tL")

	vm.Registers.A = vm.sbb(vm.Registers.L)
}

// SBB M: Subtract memory address pointed to by register pair HL from accumulator with borrow.
func (vm *CPU8080) sbb_M(data []byte) {
	vm.Logger.Debug("[9E] SBB \tM")

	vm.Registers.A = vm.sbb(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
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

// XRA A: Exclusive-OR accumulator with accumulator.
func (vm *CPU8080) xra_A(data []byte) {
	vm.Logger.Debug("[AF] XOR \tA")
	vm.xra(vm.Registers.A)
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

// ORA A: OR A with register A
func (vm *CPU8080) ora_A(data []byte) {
	vm.Logger.Debug("[B7] OR  \tA")
	vm.ora(vm.Registers.A)
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

// ORA M: OR A with memory location pointed to by register pair HL
func (vm *CPU8080) ora_M(data []byte) {
	vm.Logger.Debugf("[B6] OR  \t(HL)")
	address := toUint16(vm.Registers.H, vm.Registers.L)
	vm.ora(vm.Memory[address])
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

// CMP A: Compare A with register A
func (vm *CPU8080) cmp_A(data []byte) {
	vm.Logger.Debugf("[BF] CP  \tA")
	vm.compare(vm.Registers.A)
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

// CMP M: Compare A with memory address pointed to by register pair HL
func (vm *CPU8080) cmp_M(data []byte) {
	vm.Logger.Debugf("[BE] CP  \t(HL)")
	vm.compare(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}
