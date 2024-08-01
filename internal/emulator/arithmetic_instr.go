package emulator

// ADI D8: ADD accumulator with 8-bit immediate value.
func (vm *CPU8080) add(data []byte) {
	vm.Logger.Debugf("[C6] ADD \t$%02X", data[0])
	result := byte(uint16(vm.registers.A) + uint16(data[0]))

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.C = carryAdd(vm.registers.A, data[0])
	vm.flags.H = auxCarryAdd(vm.registers.A, data[0])
	vm.flags.setP(uint16(result))

	vm.registers.A = byte(result)
}

func inc(reg1 byte, reg2 byte) (byte, byte) {
	combined := toUint16(reg1, reg2)
	combined++

	return byte(combined >> 8), byte(combined & 0xFF)
}

// INC H: Increment register pair H.
func (vm *CPU8080) inc_HL(data []byte) {
	vm.Logger.Debugf("[23] INC \tHL")
	vm.registers.H, vm.registers.L = inc(vm.registers.H, vm.registers.L)
}

// INC D: Increment register pair D.
func (vm *CPU8080) inc_DE(data []byte) {
	vm.Logger.Debugf("[13] INC \tDE")
	vm.registers.D, vm.registers.E = inc(vm.registers.D, vm.registers.E)
}

func (vm *CPU8080) dec(data byte) byte {
	result := uint16(data) - 1

	// Handle condition bits
	vm.flags.setZ(result)
	vm.flags.setS(result)
	vm.flags.H = auxCarrySub(data, 1)
	vm.flags.setP(result)

	return byte(result)
}

// DCR B: Decrement register B.
func (vm *CPU8080) dec_B(data []byte) {
	vm.Logger.Debugf("[05] DEC \tB")
	vm.registers.B = vm.dec(vm.registers.B)
}

// DCR A: Decrement register A.
func (vm *CPU8080) dec_A(data []byte) {
	vm.Logger.Debugf("[3D] DEC \tA")
	vm.registers.A = vm.dec(vm.registers.A)
}

// DCR M: Decrement memory location pointed to by register pair HL.
func (vm *CPU8080) dec_HL(data []byte) {
	vm.Logger.Debugf("[35] DEC \t(HL)")
	memoryAddress := toUint16(vm.registers.H, vm.registers.L)
	vm.memory[memoryAddress] = vm.dec(vm.memory[memoryAddress])
}

// DCR C: Decrement register C.
func (vm *CPU8080) dec_C(data []byte) {
	vm.Logger.Debugf("[0D] DEC \tC")
	vm.registers.C = vm.dec(vm.registers.C)
}

// DAD H: Add register pair H to register pair H.
func (vm *CPU8080) dad_HL(data []byte) {
	vm.Logger.Debugf("[29] ADD \tHL,HL")
	hl := toUint16(vm.registers.H, vm.registers.L)
	doubledHL := uint32(hl) << 1

	vm.flags.C = doubledHL > 0xFFFF

	vm.registers.H = byte(doubledHL >> 8)
	vm.registers.L = byte(doubledHL)
}

// DAD D: Add register pair D to register pair H.
func (vm *CPU8080) dad_DE(data []byte) {
	vm.Logger.Debugf("[19] ADD \tHL,DE")
	de := uint32(toUint16(vm.registers.D, vm.registers.E))
	hl := uint32(toUint16(vm.registers.H, vm.registers.L))

	result := hl + de

	vm.flags.C = result > 0xFFFF

	vm.registers.H = byte(result >> 8)
	vm.registers.L = byte(result)
}

// DAD B: Add register pair B to register pair H.
func (vm *CPU8080) dad_BC(data []byte) {
	vm.Logger.Debugf("[09] ADD \tHL,BC")
	bc := uint32(toUint16(vm.registers.B, vm.registers.C))
	hl := uint32(toUint16(vm.registers.H, vm.registers.L))

	result := hl + bc

	vm.flags.C = result > 0xFFFF

	vm.registers.H = byte(result >> 8)
	vm.registers.L = byte(result)
}

// XCHG: Exchange register pairs D and H.
func (vm *CPU8080) xchg(data []byte) {
	vm.Logger.Debugf("[EB] EX  \tDE,HL")
	vm.registers.D, vm.registers.H = vm.registers.H, vm.registers.D
	vm.registers.E, vm.registers.L = vm.registers.L, vm.registers.E
}
