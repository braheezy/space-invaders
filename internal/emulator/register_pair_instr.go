package emulator

func (vm *CPU8080) push(lower, upper byte) {
	// Store value in stack, note: stack grows downwards
	vm.Memory[vm.sp-1] = upper
	vm.Memory[vm.sp-2] = lower
	vm.sp -= 2
}

// PUSH D: Push register pair D onto stack.
func (vm *CPU8080) push_DE(data []byte) {
	vm.Logger.Debugf("[D5] PUSH\tDE")
	vm.push(vm.Registers.E, vm.Registers.D)
}

// PUSH H: Push register pair H onto stack.
func (vm *CPU8080) push_HL(data []byte) {
	vm.Logger.Debugf("[E5] PUSH\tHL")
	vm.push(vm.Registers.L, vm.Registers.H)
}

// PUSH B: Push register pair B onto stack.
func (vm *CPU8080) push_BC(data []byte) {
	vm.Logger.Debugf("[C5] PUSH\tBC")
	vm.push(vm.Registers.C, vm.Registers.B)
}

// PUSH AF: Push accumulator and flags onto stack.
func (vm *CPU8080) push_AF(data []byte) {
	vm.Logger.Debugf("[F5] PUSH\tAF")
	vm.push(vm.flags.toByte(), vm.Registers.A)
}

// pop returns two bytes from the stack.
func (vm *CPU8080) pop() (byte, byte) {
	lower := vm.Memory[vm.sp]
	upper := vm.Memory[vm.sp+1]
	vm.sp += 2
	return lower, upper
}

// POP H: Pop register pair H from stack.
func (vm *CPU8080) pop_HL(data []byte) {
	vm.Logger.Debugf("[E1] POP \tHL")
	vm.Registers.L, vm.Registers.H = vm.pop()
}

// POP B: Pop register pair B from stack.
func (vm *CPU8080) pop_BC(data []byte) {
	vm.Logger.Debugf("[C1] POP \tBC")
	vm.Registers.C, vm.Registers.B = vm.pop()
}

// POP D: Pop register pair D from stack.
func (vm *CPU8080) pop_DE(data []byte) {
	vm.Logger.Debugf("[D1] POP \tDE")
	vm.Registers.E, vm.Registers.D = vm.pop()
}

// POP AF: Pop accumulator and flags from stack.
func (vm *CPU8080) pop_AF(data []byte) {
	vm.Logger.Debugf("[F1] POP \tAF")
	var fl byte
	fl, vm.Registers.A = vm.pop()
	vm.flags = *fromByte(fl)
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

// XCHG: Exchange register pairs D and H.
func (vm *CPU8080) xchg(data []byte) {
	vm.Logger.Debugf("[EB] EX  \tDE,HL")
	vm.Registers.D, vm.Registers.H = vm.Registers.H, vm.Registers.D
	vm.Registers.E, vm.Registers.L = vm.Registers.L, vm.Registers.E
}

// XTHL: Exchange top of stack with address referenced by register pair HL.
func (vm *CPU8080) xthl(data []byte) {
	vm.Logger.Debugf("[E3] EX  \t(SP),HL")
	stackL := vm.Memory[vm.sp]
	stackH := vm.Memory[vm.sp+1]

	// Exchange the values
	vm.Memory[vm.sp] = vm.Registers.L
	vm.Memory[vm.sp+1] = vm.Registers.H

	// Update the HL register pair
	vm.Registers.L = stackL
	vm.Registers.H = stackH
}

// SPHL: Load stack pointer from register pair HL.
func (vm *CPU8080) sphl(data []byte) {
	vm.Logger.Debugf("[F9] SPHL")
	vm.sp = (uint16(vm.Registers.H) << 8) | uint16(vm.Registers.L)
}
