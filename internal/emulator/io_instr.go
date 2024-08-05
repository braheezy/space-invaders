package emulator

import "fmt"

// OUT D8: Output accumulator to device at 8-bit immediate address.
func (vm *CPU8080) out(data []byte) {
	address := data[0]
	deviceName := vm.Hardware.OutDeviceName(address)
	vm.Logger.Debugf("[D3] OUT \t(%s),A", deviceName)
	vm.PC++
	err := vm.Hardware.Out(address, vm.Registers.A)
	if err != nil {
		vm.Logger.Fatal("OUT", "address", fmt.Sprintf("%04X", vm.PC-2), "error", err)
	}
}

// IN D8: Input accumulator from device at 8-bit immediate address.
func (vm *CPU8080) in(data []byte) {
	address := data[0]
	deviceName := vm.Hardware.InDeviceName(address)
	vm.Logger.Debugf("[D8] IN  \tA,(%s)", deviceName)
	vm.PC++
	result, err := vm.Hardware.In(address)
	if err != nil {
		vm.Logger.Fatal("IN", "address", fmt.Sprintf("%04X", vm.PC-2), "error", err)
	}
	vm.Registers.A = result
}
