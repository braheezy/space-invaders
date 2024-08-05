package emulator

import "testing"

func TestDecHL(t *testing.T) {
	tests := []struct {
		name                string
		initialH            byte
		initialL            byte
		memoryLocation      int
		initialMemoryValue  byte
		expectedMemoryValue byte
	}{
		{
			name:                "Normal decrement",
			initialH:            0x3A,
			initialL:            0x7C,
			memoryLocation:      0x3A7C,
			initialMemoryValue:  0x40,
			expectedMemoryValue: 0x3F,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewEmulator(&NullHardware{})
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL
			vm.Memory[tt.memoryLocation] = tt.initialMemoryValue

			vm.dcr_M([]byte{0x00, 0x00})

			if vm.Memory[tt.memoryLocation] != tt.expectedMemoryValue {
				t.Errorf("Expected memory value=0x%02X; got 0x%02X", tt.expectedMemoryValue, vm.Memory[tt.memoryLocation])
			}
		})
	}
}

func TestDcxH(t *testing.T) {
	tests := []struct {
		name      string
		initialH  byte
		initialL  byte
		expectedH byte
		expectedL byte
	}{
		{
			name:      "Normal decrement",
			initialH:  0x98,
			initialL:  0x00,
			expectedH: 0x97,
			expectedL: 0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewEmulator(&NullHardware{})
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL

			vm.dcx_H([]byte{0x00, 0x00})

			if vm.Registers.H != tt.expectedH || vm.Registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.Registers.H, vm.Registers.L)
			}
		})
	}
}
func TestIncC(t *testing.T) {
	vm := NewEmulator(&NullHardware{})
	vm.Registers.C = 0x99
	vm.inr_C(nil)
	if vm.Registers.C != 0x9A {
		t.Errorf("Expected C=0x9A, got 0x%02X", vm.Registers.C)
	}
	if vm.flags.Z {
		t.Errorf("Expected Z flag false, got true")
	}
	if !vm.flags.S {
		t.Errorf("Expected S flag false, got true")
	}
	if !vm.flags.P {
		t.Errorf("Expected P flag true, got false")
	}
	if vm.flags.C {
		t.Errorf("Expected C flag false, got true")
	}
	if vm.flags.H {
		t.Errorf("Expected H flag false, got true")
	}
}
