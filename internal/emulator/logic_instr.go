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
	vm.Logger.Debugf("[E6] AND \tA")
	vm.ana(vm.registers.A)
}
