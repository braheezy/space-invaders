package emulator

// Interrupt defines an interrupt for hardware.
type Interrupt struct {
	// Cycle is the number of cycles that should execute in the current frame before the interrupt triggers.
	Cycle int
	// Action is a function called when the interrupt triggers.
	Action func(*CPU8080)
	// Name is the name of the interrupt.
	Name string
}

// handleInterrupt performs actions when an interrupt request is received.
// The current program counter is saved to the stack and program execution changes to the
// interrupt routine address.
func (vm *CPU8080) handleInterrupt(opcode byte) {
	// Prevent other interrupts from editing CPU state at the same time.
	vm.mu.Lock()
	defer vm.mu.Unlock()

	// Check if interrupts are enabled. If not, simply return.
	if !vm.interruptsEnabled {
		return
	}

	// Disable further interrupts to prevent re-entry
	vm.interruptsEnabled = false

	// Calculate the address from the opcode (RST n: n*8)
	address := uint16((opcode - 0xC7) / 8 * 8)
	vm.Logger.Debugf("INTE $%04X-->$%04X", vm.PC, address)

	// Push the current PC onto the stack. Assumes a function exists to handle pushing words onto the stack.
	vm.push(byte(vm.PC&0xFF), byte(vm.PC>>8)&0xFF)

	// Set the PC to the ISR address.
	vm.PC = address
}
