package emulator

// ADI D8: ADD accumulator with 8-bit immediate value.
func (vm *CPU8080) adi(data []byte) {
	vm.Logger.Debugf("[C6] ADD \tA,$%02X", data[0])
	result := byte(uint16(vm.registers.A) + uint16(data[0]))

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carryAdd(vm.registers.A, data[0])
	vm.flags.H = auxCarryAdd(vm.registers.A, data[0])
	vm.flags.setP(uint16(result))

	vm.registers.A = byte(result)
	vm.pc++
}

// increment helper
func (vm *CPU8080) inc(data byte) byte {
	result := uint16(data) + 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.H = auxCarryAdd(data, 1)
	vm.flags.setP(result)

	return byte(result)
}

// INR A: Increment register A.
func (vm *CPU8080) inr_A(data []byte) {
	vm.Logger.Debugf("[3C] INR \tA")
	vm.registers.A = vm.inc(vm.registers.A)
}

// INR B: Increment register B.
func (vm *CPU8080) inr_B(data []byte) {
	vm.Logger.Debugf("[04] INR \tB")
	vm.registers.B = vm.inc(vm.registers.B)
}

// increment pair helper
func inx(reg1 byte, reg2 byte) (byte, byte) {
	combined := toUint16(reg1, reg2)
	combined++

	return byte(combined >> 8), byte(combined & 0xFF)
}

// INX H: Increment register pair H.
func (vm *CPU8080) inx_H(data []byte) {
	vm.Logger.Debugf("[23] INC \tHL")
	vm.registers.H, vm.registers.L = inx(vm.registers.H, vm.registers.L)
}

// INX D: Increment register pair D.
func (vm *CPU8080) inx_D(data []byte) {
	vm.Logger.Debugf("[13] INC \tDE")
	vm.registers.D, vm.registers.E = inx(vm.registers.D, vm.registers.E)
}

// INX B: Increment register pair B.
func (vm *CPU8080) inx_B(data []byte) {
	vm.Logger.Debugf("[03] INC \tBC")
	vm.registers.B, vm.registers.C = inx(vm.registers.B, vm.registers.C)
}

// decrement helper
func (vm *CPU8080) dcr(data byte) byte {
	result := uint16(data) - 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.C = carrySub(data, 1)
	vm.flags.H = auxCarrySub(data, 1)
	vm.flags.setP(result)

	return byte(result)
}

// DCR B: Decrement register B.
func (vm *CPU8080) dcr_B(data []byte) {
	vm.Logger.Debugf("[05] DEC \tB")
	vm.registers.B = vm.dcr(vm.registers.B)
}

// DCR A: Decrement register A.
func (vm *CPU8080) dcr_A(data []byte) {
	vm.Logger.Debugf("[3D] DEC \tA")
	vm.registers.A = vm.dcr(vm.registers.A)
}

// DCR M: Decrement memory location pointed to by register pair HL.
func (vm *CPU8080) dcr_M(data []byte) {
	vm.Logger.Debugf("[35] DEC \t(HL)")
	memoryAddress := toUint16(vm.registers.H, vm.registers.L)
	vm.memory[memoryAddress] = vm.dcr(vm.memory[memoryAddress])
}

// DCR C: Decrement register C.
func (vm *CPU8080) dcr_C(data []byte) {
	vm.Logger.Debugf("[0D] DEC \tC")
	vm.registers.C = vm.dcr(vm.registers.C)
}

// decrement pair helper
func decPair(reg1 byte, reg2 byte) (byte, byte) {
	combined := toUint16(reg1, reg2)
	combined--

	return byte(combined >> 8), byte(combined & 0xFF)
}

// DCX H: Decrement register pair H.
func (vm *CPU8080) dcx_H(data []byte) {
	vm.Logger.Debugf("[2B] DEC \tHL")
	vm.registers.H, vm.registers.L = decPair(vm.registers.H, vm.registers.L)
}

// SUI D8: Subtract immediate value from accumulator.
func (vm *CPU8080) sui(data []byte) {
	vm.Logger.Debugf("[D6] SUB \t$%02X", data[0])
	result := uint16(vm.registers.A) - uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carrySub(vm.registers.A, data[0])
	vm.flags.H = auxCarrySub(vm.registers.A, data[0])
	vm.flags.setP(uint16(result))

	vm.registers.A = byte(result)
	vm.pc++
}

// DAD H: Add register pair H to register pair H.
func (vm *CPU8080) dad_H(data []byte) {
	vm.Logger.Debugf("[29] ADD \tHL,HL")
	hl := toUint16(vm.registers.H, vm.registers.L)
	doubledHL := uint32(hl) << 1

	vm.flags.C = doubledHL > 0xFFFF

	vm.registers.H = byte(doubledHL >> 8)
	vm.registers.L = byte(doubledHL)
}

// DAD D: Add register pair D to register pair H.
func (vm *CPU8080) dad_D(data []byte) {
	vm.Logger.Debugf("[19] ADD \tHL,DE")
	de := uint32(toUint16(vm.registers.D, vm.registers.E))
	hl := uint32(toUint16(vm.registers.H, vm.registers.L))

	result := hl + de

	vm.flags.C = result > 0xFFFF

	vm.registers.H = byte(result >> 8)
	vm.registers.L = byte(result)
}

// DAD B: Add register pair B to register pair H.
func (vm *CPU8080) dad_B(data []byte) {
	vm.Logger.Debugf("[09] ADD \tHL,BC")
	bc := uint32(toUint16(vm.registers.B, vm.registers.C))
	hl := uint32(toUint16(vm.registers.H, vm.registers.L))

	result := hl + bc

	vm.flags.C = result > 0xFFFF

	vm.registers.H = byte(result >> 8)
	vm.registers.L = byte(result)
}
