package emulator

import "testing"

func TestSTAX(t *testing.T) {

	vm := NewCPU8080(&[]byte{}, nil)
	vm.Registers.A = 0x42
	vm.Registers.B = 0x3F
	vm.Registers.C = 0x16
	vm.stax_B(nil)

	if vm.Memory[0x3F16] != 0x42 {
		t.Error("LDAX failed")
	}
}
