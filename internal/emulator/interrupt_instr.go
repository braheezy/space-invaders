package emulator

// EI: Enable interrupts.
func (vm *CPU8080) ei(data []byte) {
	vm.Logger.Debugf("[FB] EI")
	vm.interruptsEnabled = true
}

// DI: Disable interrupts.
func (vm *CPU8080) di(data []byte) {
	vm.Logger.Debugf("[F3] DI")
	vm.interruptsEnabled = false
}
