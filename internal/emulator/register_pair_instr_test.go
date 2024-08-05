package emulator

import "testing"

func TestInxHL(t *testing.T) {
	tests := []struct {
		name          string
		initialH      byte
		initialL      byte
		expectedH     byte
		expectedL     byte
		carryFlagSet  bool
		expectedCarry bool
	}{
		{
			name:          "Normal increment",
			initialH:      0x38,
			initialL:      0xFF,
			expectedH:     0x39,
			expectedL:     0x00,
			carryFlagSet:  false,
			expectedCarry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewEmulator(&NullHardware{})
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL
			vm.flags.C = tt.carryFlagSet

			vm.inx_H([]byte{0x00, 0x00})

			if vm.Registers.H != tt.expectedH || vm.Registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.Registers.H, vm.Registers.L)
			}
		})
	}
}

func TestDadDE(t *testing.T) {
	tests := []struct {
		name          string
		initialD      byte
		initialE      byte
		initialH      byte
		initialL      byte
		expectedH     byte
		expectedL     byte
		expectedCarry bool
	}{
		{
			name:          "Normal increment",
			initialD:      0x33,
			initialE:      0x9F,
			initialH:      0xA1,
			initialL:      0x7B,
			expectedH:     0xD5,
			expectedL:     0x1A,
			expectedCarry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewEmulator(&NullHardware{})
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL
			vm.Registers.D = tt.initialD
			vm.Registers.E = tt.initialE
			vm.flags.C = true

			vm.dad_D([]byte{0x00, 0x00})

			if vm.Registers.H != tt.expectedH || vm.Registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.Registers.H, vm.Registers.L)
			}
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("expected C flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}

func TestXCHG(t *testing.T) {

	vm := NewEmulator(&NullHardware{})
	vm.Registers.D = 0x33
	vm.Registers.E = 0x55
	vm.Registers.H = 0x00
	vm.Registers.L = 0xFF

	vm.xchg(nil)

	if vm.Registers.D != 0x00 {
		t.Errorf("Expected D to be 0x00, got %02X", vm.Registers.D)
	}
	// E should be 0xFF
	if vm.Registers.E != 0xFF {
		t.Errorf("Expected E to be 0xFF, got %02X", vm.Registers.E)
	}
	// H should be 0x33
	if vm.Registers.H != 0x33 {
		t.Errorf("Expected H to be 0x33, got %02X", vm.Registers.H)
	}
	// L should be 0x55
	if vm.Registers.L != 0x55 {
		t.Errorf("Expected L to be 0x55, got %02X", vm.Registers.L)
	}

}

func TestDAD_D(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewEmulator(&NullHardware{})

	// Test cases
	tests := []struct {
		name          string
		initialD      byte
		initialE      byte
		initialH      byte
		initialL      byte
		expectedH     byte
		expectedL     byte
		expectedCarry bool
	}{
		{
			name:          "Double without carry",
			initialD:      0x33,
			initialE:      0x9F,
			initialH:      0xA1,
			initialL:      0x7B,
			expectedH:     0xD5,
			expectedL:     0x1A,
			expectedCarry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the initial D and E register values
			vm.Registers.D = tt.initialD
			vm.Registers.E = tt.initialE
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL

			vm.dad_D(nil)

			// Check if the H and L registers have the expected values
			if vm.Registers.H != tt.expectedH || vm.Registers.L != tt.expectedL {
				t.Errorf("Expected HL value %02X%02X, got %02X%02X", tt.expectedH, tt.expectedL, vm.Registers.H, vm.Registers.L)
			}

			// Check if the carry flag is set as expected
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected carry flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}
