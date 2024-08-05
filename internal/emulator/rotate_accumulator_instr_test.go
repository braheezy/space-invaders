package emulator

import "testing"

func TestRLC(t *testing.T) {

	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewEmulator(&NullHardware{})

	tests := []struct {
		name      string
		initialA  byte
		expectedA byte
		expectedC bool
	}{
		{
			name:      "RLC with zero",
			initialA:  0x00,
			expectedA: 0x00,
			expectedC: false,
		},
		{
			name:      "RLC with non-zero",
			initialA:  0x0F,
			expectedA: 0x1E,
			expectedC: false,
		},
		{
			name:      "RLC with carry",
			initialA:  0x0F2,
			expectedA: 0x0E5,
			expectedC: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and set the initial accumulator value
			vm.flags = flags{}
			vm.Registers.A = tt.initialA

			// Execute the RLC function
			vm.rlc([]byte{0})

			// Check the accumulator's value and flags
			if vm.Registers.A != tt.expectedA {
				t.Errorf("%s: expected accumulator %02X, got %02X", tt.name, tt.expectedA, vm.Registers.A)
			}
			if vm.flags.C != tt.expectedC {
				t.Errorf("%s: expected C flag %t, got %t", tt.name, tt.expectedC, vm.flags.C)
			}
		})
	}
}

func TestRAR(t *testing.T) {
	vm := NewEmulator(&NullHardware{})
	vm.Registers.A = 0x6A
	vm.flags.C = true

	vm.rar(nil)

	expectedA := byte(0xB5)
	expectedC := false

	if vm.Registers.A != expectedA {
		t.Errorf("expected accumulator %02X, got %02X", expectedA, vm.Registers.A)
	}
	if vm.flags.C != expectedC {
		t.Errorf("expected C flag %t, got %t", expectedC, vm.flags.C)
	}

}

func TestRAL(t *testing.T) {

	vm := NewEmulator(&NullHardware{})
	vm.Registers.A = 0xB5

	vm.ral(nil)

	if vm.Registers.A != 0x6A {
		t.Errorf("Expected accumulator to be 0x6A, got %02X", vm.Registers.A)
	}
	if !vm.flags.C {
		t.Error("Expected carry flag to be set")
	}
}

// TestRRC tests the rrc function for rotating the accumulator right
// and updating the carry flag.
func TestRRC(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewEmulator(&NullHardware{})

	// Test cases
	tests := []struct {
		name          string
		initialA      byte
		expectedA     byte
		expectedCarry bool
	}{
		{
			name:          "Rotate with carry set",
			initialA:      0b00000001,
			expectedA:     0b10000000,
			expectedCarry: true,
		},
		{
			name:          "Rotate with carry unset",
			initialA:      0b00000010,
			expectedA:     0b00000001,
			expectedCarry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the accumulator to the initial value
			vm.Registers.A = tt.initialA

			// Execute the rrc function
			vm.rrc(nil)

			// Check if the accumulator has the expected value
			if vm.Registers.A != tt.expectedA {
				t.Errorf("Expected accumulator value %08b, got %08b", tt.expectedA, vm.Registers.A)
			}

			// Check if the carry flag is set as expected
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected carry flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}
