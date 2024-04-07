package emulator

func (vm *CPU8080) push(lower, upper byte) {
	// Store value in stack, note: stack grows downwards
	vm.memory[vm.sp-1] = upper
	vm.memory[vm.sp-2] = lower
	vm.sp -= 2
}

// PUSH D: Push register pair D onto stack.
func (vm *CPU8080) push_DE(data []byte) {
	vm.Logger.Debugf("[D5] PUSH\tDE")
	vm.push(vm.registers.E, vm.registers.D)
}

// PUSH H: Push register pair H onto stack.
func (vm *CPU8080) push_HL(data []byte) {
	vm.Logger.Debugf("[E5] PUSH\tHL")
	vm.push(vm.registers.L, vm.registers.H)
}

// PUSH B: Push register pair B onto stack.
func (vm *CPU8080) push_BC(data []byte) {
	vm.Logger.Debugf("[C5] PUSH\tBC")
	vm.push(vm.registers.C, vm.registers.B)
}

// PUSH AF: Push accumulator and flags onto stack.
func (vm *CPU8080) push_AF(data []byte) {
	vm.Logger.Debugf("[F5] PUSH\tAF")
	vm.push(vm.flags.toByte(), vm.registers.A)
}

// pop returns two bytes from the stack.
func (vm *CPU8080) pop() (byte, byte) {
	lower := vm.memory[vm.sp]
	upper := vm.memory[vm.sp+1]
	vm.sp += 2
	return lower, upper
}

// POP H: Pop register pair H from stack.
func (vm *CPU8080) pop_HL(data []byte) {
	vm.Logger.Debugf("[E1] POP \tHL")
	vm.registers.L, vm.registers.H = vm.pop()
}

// POP B: Pop register pair B from stack.
func (vm *CPU8080) pop_BC(data []byte) {
	vm.Logger.Debugf("[C1] POP \tBC")
	vm.registers.C, vm.registers.B = vm.pop()
}

// POP D: Pop register pair D from stack.
func (vm *CPU8080) pop_DE(data []byte) {
	vm.Logger.Debugf("[D1] POP \tDE")
	vm.registers.E, vm.registers.D = vm.pop()
}

// POP AF: Pop accumulator and flags from stack.
func (vm *CPU8080) pop_AF(data []byte) {
	vm.Logger.Debugf("[F1] POP \tAF")
	var fl byte
	fl, vm.registers.A = vm.pop()
	vm.flags = *fromByte(fl)
}
