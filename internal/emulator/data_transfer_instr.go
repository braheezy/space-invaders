package emulator

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) load_SP(data []byte) {
	operand := toUint16(data[1], data[0])
	vm.Logger.Debugf("[31] LD  \tSP,$%04X", operand)
	vm.sp = operand
	vm.pc += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) load_BC(data []byte) {
	vm.Logger.Debugf("[01] LD  \tB,$%04X", toUint16(data[1], data[0]))
	vm.registers.C = data[0]
	vm.registers.B = data[1]
	vm.pc += 2
}

// MVI A, D8: Move 8-bit immediate value into accumulator.
func (vm *CPU8080) moveI_A(data []byte) {
	vm.Logger.Debugf("[3E] LD  \tA,$%02X", data[0])
	vm.registers.A = data[0]
	vm.pc++
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveI_B(data []byte) {
	vm.Logger.Debugf("[06] LD  \tB,$%02X", data[0])
	vm.registers.B = data[0]
	vm.pc++
}

// MVI C, D8: Move 8-bit immediate value into register C.
func (vm *CPU8080) moveI_C(data []byte) {
	vm.Logger.Debugf("[0E] LD  \tC,$%02X", data[0])
	vm.registers.C = data[0]
	vm.pc++
}

// MVI H, D8: Move 8-bit immediate value into register H.
func (vm *CPU8080) moveI_H(data []byte) {
	vm.Logger.Debugf("[26] LD  \tH,$%02X", data[0])
	vm.registers.H = data[0]
	vm.pc++
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) load_DE(data []byte) {
	vm.Logger.Debugf("[11] LD  \tDE,$%04X", toUint16(data[1], data[0]))
	vm.registers.E = data[0]
	vm.registers.D = data[1]
	vm.pc += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) load_HL(data []byte) {
	vm.Logger.Debugf("[21] LD  \tHL,$%04X", toUint16(data[1], data[0]))
	vm.registers.L = data[0]
	vm.registers.H = data[1]
	vm.pc += 2
}

// LDAX D: Load value from address in register pair D into accumulator.
func (vm *CPU8080) load_DEA(data []byte) {
	address := toUint16(vm.registers.D, vm.registers.E)
	vm.Logger.Debugf("[1A] LD  \tA,(DE)")
	vm.registers.A = vm.memory[address]
}

// MOV M, A: Move value from accumulator into register pair H.
func (vm *CPU8080) load_HLA(data []byte) {
	address := toUint16(vm.registers.H, vm.registers.L)
	vm.Logger.Debugf("[77] LD  \t(HL),A ($%04X)", address)
	vm.memory[address] = vm.registers.A
}

// MOV L,A: Load value from accumulator into register L.
func (vm *CPU8080) move_AL(data []byte) {
	vm.Logger.Debugf("[6F] LD  \tL,A")
	vm.registers.A = vm.registers.L
}

// MVI HL: Move 8-bit immediate value into memory address from register pair HL
func (vm *CPU8080) moveI_HL(data []byte) {
	address := toUint16(vm.registers.H, vm.registers.L)
	vm.Logger.Debugf("[36] LD  \t(HL),$%02X", data[0])
	vm.memory[address] = data[0]
	vm.pc++
}

// MOV E, HL: Move memory location pointed to by register pair HL into register E.
func (vm *CPU8080) moveHL_E(data []byte) {
	vm.Logger.Debugf("[5E] LD  \tE,(HL)")
	vm.registers.E = vm.memory[toUint16(vm.registers.H, vm.registers.L)]
}

// MOV D, HL: Move memory location pointed to by register pair HL into register D.
func (vm *CPU8080) moveHL_D(data []byte) {
	vm.Logger.Debugf("[56] LD  \tD,(HL)")
	vm.registers.D = vm.memory[toUint16(vm.registers.H, vm.registers.L)]
}

// MOV A, HL: Move memory location pointed to by register pair HL into register A.
func (vm *CPU8080) moveHL_A(data []byte) {
	vm.Logger.Debugf("[7E] LD  \tA,(HL)")
	vm.registers.A = vm.memory[toUint16(vm.registers.H, vm.registers.L)]
}

// MOV H, HL: Move memory location pointed to by register pair HL into register H.
func (vm *CPU8080) moveHL_H(data []byte) {
	vm.Logger.Debugf("[66] LD  \tH,(HL)")
	vm.registers.H = vm.memory[toUint16(vm.registers.H, vm.registers.L)]
}

// MOV A,H: Move value from register H into accumulator.
func (vm *CPU8080) move_HA(data []byte) {
	vm.Logger.Debugf("[7E] LD  \tA,H")
	vm.registers.A = vm.registers.H
}

// MOV A,H: Move value from register D into accumulator.
func (vm *CPU8080) move_DA(data []byte) {
	vm.Logger.Debugf("[7C] LD  \tA,D")
	vm.registers.A = vm.registers.D
}

// MOV A,E: Move value from register E into accumulator.
func (vm *CPU8080) move_EA(data []byte) {
	vm.Logger.Debugf("[7B] LD  \tA,E")
	vm.registers.A = vm.registers.E
}

// STA A16: Store accumulator in 16-bit immediate address.
func (vm *CPU8080) store_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[32] LD  \t$%04X,A", address)
	vm.memory[address] = vm.registers.A
	vm.pc += 2
}

// LDA A16: Load accumulator from 16-bit immediate address.
func (vm *CPU8080) load_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[3A] LD  \tA,$%04X", address)
	vm.registers.A = vm.memory[address]
	vm.pc += 2
}
