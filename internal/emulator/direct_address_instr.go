package emulator

// SHLD A16: Store register pair HL into 16-bit immediate address.
func (vm *CPU8080) store_HL(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[22] LD  \t$%04X,HL", address)
	vm.Memory[address] = vm.Registers.L
	vm.Memory[address+1] = vm.Registers.H
	vm.PC += 2
}

// LHLD A16: Load register pair HL from 16-bit immediate address.
func (vm *CPU8080) loadImm_HL(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[2A] LD  \tHL,$%04X", address)
	vm.Registers.L = vm.Memory[address]
	vm.Registers.H = vm.Memory[address+1]
	vm.PC += 2
}

// STA A16: Store accumulator in 16-bit immediate address.
func (vm *CPU8080) store_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[32] LD  \t$%04X,A", address)
	vm.Memory[address] = vm.Registers.A
	vm.PC += 2
}

// LDA A16: Load accumulator from 16-bit immediate address.
func (vm *CPU8080) load_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[3A] LD  \tA,$%04X", address)
	vm.Registers.A = vm.Memory[address]
	vm.PC += 2
}
