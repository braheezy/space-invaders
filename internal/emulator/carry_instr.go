package emulator

// STC: Set carry bit to 1
func (vm *CPU8080) set_C(data []byte) {
	vm.Logger.Debug("STC")
	vm.flags.C = true
}

// CMC: Complement carry bit
func (vm *CPU8080) cmc(data []byte) {
	vm.Logger.Debug("CMC")
	vm.flags.C = !vm.flags.C
}
