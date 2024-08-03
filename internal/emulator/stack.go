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
