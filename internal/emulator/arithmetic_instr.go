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

// ADI: ADD accumulator with 8-bit immediate value.
func (vm *CPU8080) adi(data []byte) {
	vm.Logger.Debugf("[C6] ADD \tA,$%02X", data[0])

	vm.Registers.A = vm.add(data[0])
	vm.PC++
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

// ACI: ADD accumulator with 8-bit immediate value with carry.
func (vm *CPU8080) aci(data []byte) {
	carry := byte(0)
	if vm.flags.C {
		carry = 1
	}
	vm.Registers.A = vm.add(data[0] + carry)

	vm.PC++
}

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

// INX SP: Increment stack pointer.
func (vm *CPU8080) inx_SP(data []byte) {
	vm.Logger.Debugf("[33] INC \tSP")
	vm.sp++
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

// DAD B: Add register pair B to register pair H.

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

// SUI: Subtract immediate value from accumulator.
func (vm *CPU8080) sui(data []byte) {
	vm.Logger.Debugf("[D6] SUB \t$%02X", data[0])

	vm.Registers.A = vm.sub(data[0])
	vm.PC++
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
	vm.Logger.Debug("[09] ADD \tHL,BC")
	bc := uint32(toUint16(vm.Registers.B, vm.Registers.C))
	hl := uint32(toUint16(vm.Registers.H, vm.Registers.L))

	result := hl + bc

	vm.flags.C = result > 0xFFFF

	vm.Registers.H = byte(result >> 8)
	vm.Registers.L = byte(result)
}

// DAD SP: Add stack pointer to register pair H.
func (vm *CPU8080) dad_SP(data []byte) {
	vm.Logger.Debug("[39] ADD \tHL,SP")
	hl := uint32(toUint16(vm.Registers.H, vm.Registers.L))

	result := hl + uint32(vm.sp)

	vm.flags.C = result > 0xFFFF

	vm.Registers.H = byte(result >> 8)
	vm.Registers.L = byte(result)
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

// ADC A: Add accumulator with register A and carry.
func (vm *CPU8080) adc_A(data []byte) {
	vm.Logger.Debug("[8F] ADC \tA")

	vm.Registers.A = vm.adc(vm.Registers.A)
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

// SBB A: Subtract register A from accumulator with borrow.
func (vm *CPU8080) sbb_A(data []byte) {
	vm.Logger.Debug("[9F] SBB \tA")

	vm.Registers.A = vm.sbb(vm.Registers.A)
}
