package emulator

// STC: Set carry bit to 1
func (vm *CPU8080) set_C(data []byte) {
	vm.flags.C = true
}
