package emulator

// add helper
func (vm *CPU8080) add(data byte) byte {
	result := uint16(vm.Registers.A) + uint16(data)

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carryAdd(vm.Registers.A, data)
	vm.flags.H = auxCarryAdd(vm.Registers.A, data)
	vm.flags.setP(uint16(result))

	return byte(result)
}

// ADI: ADD accumulator with 8-bit immediate value.
func (vm *CPU8080) adi(data []byte) {
	vm.Logger.Debugf("[C6] ADD \tA,$%02X", data[0])

	vm.Registers.A = vm.add(data[0])
	vm.PC++
}

// ADD B: ADD accumulator with register B.
func (vm *CPU8080) add_B(data []byte) {
	vm.Logger.Debug("[80] ADD \tA,B")

	vm.Registers.A = vm.add(vm.Registers.B)
}

// ADD E: ADD accumulator with register E.
func (vm *CPU8080) add_E(data []byte) {
	vm.Logger.Debug("[83] ADD \tA,E")

	vm.Registers.A = vm.add(vm.Registers.E)
}

// ADD M: ADD accumulator with memory address pointed to by register pair HL
func (vm *CPU8080) add_M(data []byte) {
	vm.Logger.Debug("[86] ADD \tA,(HL)")

	vm.Registers.A = vm.add(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
}

// ADD L: ADD accumulator with register L.
func (vm *CPU8080) add_L(data []byte) {
	vm.Logger.Debug("[85] ADD \tA,L")

	vm.Registers.A = vm.add(vm.Registers.L)
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
	vm.Logger.Debugf("[3C] INC \tA")
	vm.Registers.A = vm.inc(vm.Registers.A)
}

// INR B: Increment register B.
func (vm *CPU8080) inr_B(data []byte) {
	vm.Logger.Debugf("[04] INC \tB")
	vm.Registers.B = vm.inc(vm.Registers.B)
}

// INR C: Increment register C.
func (vm *CPU8080) inr_C(data []byte) {
	vm.Logger.Debugf("[0C] INC \tC")
	vm.Registers.C = vm.inc(vm.Registers.C)
}

// INR D: Increment register D.
func (vm *CPU8080) inr_D(data []byte) {
	vm.Logger.Debugf("[14] INC \tD")
	vm.Registers.D = vm.inc(vm.Registers.D)
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
	vm.Registers.H, vm.Registers.L = inx(vm.Registers.H, vm.Registers.L)
}

// INX D: Increment register pair D.
func (vm *CPU8080) inx_D(data []byte) {
	vm.Logger.Debugf("[13] INC \tDE")
	vm.Registers.D, vm.Registers.E = inx(vm.Registers.D, vm.Registers.E)
}

// INX B: Increment register pair B.
func (vm *CPU8080) inx_B(data []byte) {
	vm.Logger.Debugf("[03] INC \tBC")
	vm.Registers.B, vm.Registers.C = inx(vm.Registers.B, vm.Registers.C)
}

// decrement helper
func (vm *CPU8080) dcr(data byte) byte {
	result := uint16(data) - 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.H = auxCarrySub(data, 1)
	vm.flags.setP(result)

	return byte(result)
}

// DCR A: Decrement register A.
func (vm *CPU8080) dcr_A(data []byte) {
	vm.Logger.Debugf("[3D] DEC \tA")
	vm.Registers.A = vm.dcr(vm.Registers.A)
}

// DCR B: Decrement register B.
func (vm *CPU8080) dcr_B(data []byte) {
	vm.Logger.Debugf("[05] DEC \tB")
	vm.Registers.B = vm.dcr(vm.Registers.B)
}

// DCR C: Decrement register C.
func (vm *CPU8080) dcr_C(data []byte) {
	vm.Logger.Debugf("[0D] DEC \tC")
	vm.Registers.C = vm.dcr(vm.Registers.C)
}

// DCR D: Decrement register D.
func (vm *CPU8080) dcr_D(data []byte) {
	vm.Logger.Debugf("[15] DEC \tD")
	vm.Registers.D = vm.dcr(vm.Registers.D)
}

// DCR M: Decrement memory location pointed to by register pair HL.
func (vm *CPU8080) dcr_M(data []byte) {
	vm.Logger.Debugf("[35] DEC \t(HL)")
	memoryAddress := toUint16(vm.Registers.H, vm.Registers.L)
	vm.Memory[memoryAddress] = vm.dcr(vm.Memory[memoryAddress])
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
	vm.Registers.H, vm.Registers.L = decPair(vm.Registers.H, vm.Registers.L)
}

// SUI: Subtract immediate value from accumulator.
func (vm *CPU8080) sui(data []byte) {
	vm.Logger.Debugf("[D6] SUB \t$%02X", data[0])
	result := uint16(vm.Registers.A) - uint16(data[0])

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carrySub(vm.Registers.A, data[0])
	vm.flags.H = auxCarrySub(vm.Registers.A, data[0])
	vm.flags.setP(uint16(result))

	vm.Registers.A = byte(result)
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
	// vm.flags.H = (vm.registers.A & 0xF) < ((data[0] + carry) & 0xF)

	vm.Registers.A = byte(result)
	vm.PC++
}

// DAD H: Add register pair H to register pair H.
func (vm *CPU8080) dad_H(data []byte) {
	vm.Logger.Debugf("[29] ADD \tHL,HL")
	hl := toUint16(vm.Registers.H, vm.Registers.L)
	doubledHL := uint32(hl) << 1

	vm.flags.C = doubledHL > 0xFFFF

	vm.Registers.H = byte(doubledHL >> 8)
	vm.Registers.L = byte(doubledHL)
}

// DAD D: Add register pair D to register pair H.
func (vm *CPU8080) dad_D(data []byte) {
	vm.Logger.Debugf("[19] ADD \tHL,DE")
	de := uint32(toUint16(vm.Registers.D, vm.Registers.E))
	hl := uint32(toUint16(vm.Registers.H, vm.Registers.L))

	result := hl + de

	vm.flags.C = result > 0xFFFF

	vm.Registers.H = byte(result >> 8)
	vm.Registers.L = byte(result)
}

// DAD B: Add register pair B to register pair H.
func (vm *CPU8080) dad_B(data []byte) {
	vm.Logger.Debugf("[09] ADD \tHL,BC")
	bc := uint32(toUint16(vm.Registers.B, vm.Registers.C))
	hl := uint32(toUint16(vm.Registers.H, vm.Registers.L))

	result := hl + bc

	vm.flags.C = result > 0xFFFF

	vm.Registers.H = byte(result >> 8)
	vm.Registers.L = byte(result)
}
