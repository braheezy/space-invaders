package emulator

// increment helper
func (vm *CPU8080) inc(data byte) byte {
	result := data + 1

	// Handle condition bits
	vm.flags.setZ(uint16(result))
	vm.flags.setS(uint16(result))
	vm.flags.H = auxCarryAdd(data, 1)
	vm.flags.setP(uint16(result))

	return result
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

// INR E: Increment register E.
func (vm *CPU8080) inr_E(data []byte) {
	vm.Logger.Debugf("[1C] INC \tE")
	vm.Registers.E = vm.inc(vm.Registers.E)
}

// INR H: Increment register H.
func (vm *CPU8080) inr_H(data []byte) {
	vm.Logger.Debugf("[24] INC \tH")
	vm.Registers.H = vm.inc(vm.Registers.H)
}

// INR L: Increment register L.
func (vm *CPU8080) inr_L(data []byte) {
	vm.Logger.Debugf("[2C] INC \tL")
	vm.Registers.L = vm.inc(vm.Registers.L)
}

// INR M: Increment memory address pointed to by register pair HL.
func (vm *CPU8080) inr_M(data []byte) {
	vm.Logger.Debugf("[34] INC \tM")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.inc(vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)])
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

// DCR E: Decrement register E.
func (vm *CPU8080) dcr_E(data []byte) {
	vm.Logger.Debugf("[1D] DEC \tE")
	vm.Registers.E = vm.dcr(vm.Registers.E)
}

// DCR H: Decrement register H.
func (vm *CPU8080) dcr_H(data []byte) {
	vm.Logger.Debugf("[25] DEC \tH")
	vm.Registers.H = vm.dcr(vm.Registers.H)
}

// DCR L: Decrement register L.
func (vm *CPU8080) dcr_L(data []byte) {
	vm.Logger.Debugf("[2D] DEC \tL")
	vm.Registers.L = vm.dcr(vm.Registers.L)
}

// DCR M: Decrement memory location pointed to by register pair HL.
func (vm *CPU8080) dcr_M(data []byte) {
	vm.Logger.Debugf("[35] DEC \t(HL)")
	memoryAddress := toUint16(vm.Registers.H, vm.Registers.L)
	vm.Memory[memoryAddress] = vm.dcr(vm.Memory[memoryAddress])
}

// CMA: Complement accumulator.
func (vm *CPU8080) cma(data []byte) {
	vm.Logger.Debugf("[2F] CMA")
	vm.Registers.A = ^vm.Registers.A
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

// decrement pair helper
func decPair(reg1 byte, reg2 byte) (byte, byte) {
	combined := toUint16(reg1, reg2)
	combined--

	return byte(combined >> 8), byte(combined & 0xFF)
}

// DCX H: Decrement register pair H.
func (vm *CPU8080) dcx_H(data []byte) {
	vm.Logger.Debug("[2B] DEC \tHL")
	vm.Registers.H, vm.Registers.L = decPair(vm.Registers.H, vm.Registers.L)
}

// DCX B: Decrement register pair B.
func (vm *CPU8080) dcx_B(data []byte) {
	vm.Logger.Debug("[0B] DEC \tBC")
	vm.Registers.B, vm.Registers.C = decPair(vm.Registers.B, vm.Registers.C)
}

// DCX D: Decrement register pair D.
func (vm *CPU8080) dcx_D(data []byte) {
	vm.Logger.Debug("[1B] DEC \tDE")
	vm.Registers.D, vm.Registers.E = decPair(vm.Registers.D, vm.Registers.E)
}

// DCX SP: Decrement stack pointer
func (vm *CPU8080) dcx_SP(data []byte) {
	vm.Logger.Debug("[3B] DEC \tSP")
	vm.sp--
}
