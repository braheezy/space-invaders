package emulator

type InterruptCondition struct {
	Cycle  int
	Action func(*CPU8080)
	Name   string
}

func (vm *CPU8080) handleInterrupt(opcode byte) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Check if interrupts are enabled. If not, simply return.
	if !vm.InterruptsEnabled {
		return
	}

	// Disable further interrupts to prevent re-entry
	vm.InterruptsEnabled = false

	// Calculate the address from the opcode (RST n: n*8)
	address := uint16((opcode - 0xC7) / 8 * 8)
	vm.Logger.Debugf("INTE $%04X-->$%04X", vm.pc, address)

	// Push the current PC onto the stack. Assumes a function exists to handle pushing words onto the stack.
	vm.push(byte(vm.pc&0xFF), byte(vm.pc>>8)&0xFF)

	// Set the PC to the ISR address.
	vm.pc = address
}

// EI: Enable interrupts.
func (vm *CPU8080) ei(data []byte) {
	vm.Logger.Debugf("[FB] EI")
	vm.InterruptsEnabled = true
}

// DI: Disable interrupts.
func (vm *CPU8080) di(data []byte) {
	vm.Logger.Debugf("[F3] DI")
	vm.InterruptsEnabled = false
}
