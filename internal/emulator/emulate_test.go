package emulator

import "testing"

// TestRRC tests the rrc function for rotating the accumulator right
// and updating the carry flag.
func TestRRC(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	// Test cases
	tests := []struct {
		name          string
		initialA      byte // Initial accumulator value
		expectedA     byte // Expected accumulator value after rotation
		expectedCarry bool // Expected carry flag state
	}{
		{
			name:          "Rotate with carry set",
			initialA:      0b00000001, // Binary for clarity
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
			vm.registers.A = tt.initialA

			// Execute the rrc function
			vm.rrc(nil) // The data slice is not used in the function

			// Check if the accumulator has the expected value
			if vm.registers.A != tt.expectedA {
				t.Errorf("Expected accumulator value %08b, got %08b", tt.expectedA, vm.registers.A)
			}

			// Check if the carry flag is set as expected
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected carry flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}

// TestDAD_DE tests the dad_DE function for doubling the DE register pair value
// and storing the result in the HL register pair, with the correct carry flag update.
func TestDAD_DE(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	// Test cases
	tests := []struct {
		name          string
		initialD      byte // Initial D register value
		initialE      byte // Initial E register value
		expectedH     byte // Expected H register value after doubling DE
		expectedL     byte // Expected L register value after doubling DE
		expectedCarry bool // Expected carry flag state
	}{
		{
			name:          "Double without carry",
			initialD:      0x12, // Example values
			initialE:      0x34,
			expectedH:     0x24, // Result of doubling 0x1234
			expectedL:     0x68,
			expectedCarry: false,
		},
		{
			name:          "Double with carry",
			initialD:      0xFF, // Example values near the limit
			initialE:      0xFF,
			expectedH:     0xFF,
			expectedL:     0xFE,
			expectedCarry: true, // Because doubling exceeds 0xFFFF
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the initial D and E register values
			vm.registers.D = tt.initialD
			vm.registers.E = tt.initialE

			// Execute the dad_DE function
			vm.dad_D(nil) // The data slice is not used in the function

			// Check if the H and L registers have the expected values
			if vm.registers.H != tt.expectedH || vm.registers.L != tt.expectedL {
				t.Errorf("Expected HL value %02X%02X, got %02X%02X", tt.expectedH, tt.expectedL, vm.registers.H, vm.registers.L)
			}

			// Check if the carry flag is set as expected
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected carry flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}

// TestCMP tests the cmp function which compares an 8-bit immediate value with the accumulator.
func TestCMP_B(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	// Test cases
	tests := []struct {
		name               string
		initialA           byte // Initial accumulator value
		initialB           byte // Immediate value to compare with
		initialC, initialZ bool
		expectedZero       bool // Expected Zero flag
		expectedSign       bool // Expected Sign flag
		expectedCarry      bool // Expected Carry flag
		expectedAuxCarry   bool // Expected Auxiliary Carry flag
		expectedParity     bool // Expected Parity flag
	}{
		{
			name:             "Accumulator greater",
			initialA:         0x0F,
			initialB:         0x07,
			expectedZero:     false,
			expectedSign:     false, // Result is positive
			expectedCarry:    false, // No borrow
			expectedAuxCarry: false, //  No aux borrow
			expectedParity:   false, // Odd parity of result (0x08)
		},
		{
			name:             "Accumulator equal",
			initialA:         0x12,
			initialB:         0x12,
			expectedZero:     true,
			expectedSign:     false, // Zero result, sign flag is 0
			expectedCarry:    false, // Equal, no borrow
			expectedAuxCarry: false, // No borrow in any nibble
			expectedParity:   true,  // Even parity (result is 0)
		},
		{
			name:             "docs 2",
			initialA:         0x02,
			initialB:         0x05,
			expectedZero:     false,
			expectedSign:     true,  // Negative result (considering 8-bit unsigned overflow)
			expectedCarry:    true,  // Borrow occurs
			expectedAuxCarry: true,  // Aux carry occurs if there's a borrow from bit 4 to bit 3
			expectedParity:   false, // Odd parity of result (0xFF)
			initialC:         false,
			initialZ:         true,
		},
		{
			name:             "docs 1",
			initialA:         0x0A,
			initialB:         0x05,
			expectedZero:     false,
			expectedSign:     false,
			expectedCarry:    false,
			expectedAuxCarry: false,
			expectedParity:   true,
			initialC:         true,
			initialZ:         true,
		},
		{
			name:             "docs 3",
			initialA:         0xEB,
			initialB:         0x05,
			expectedZero:     false,
			expectedSign:     true,
			expectedCarry:    false,
			expectedAuxCarry: false,
			expectedParity:   false,
			initialC:         true,
			initialZ:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the initial A register value
			vm.registers.A = tt.initialA
			vm.registers.B = tt.initialB
			vm.flags.C = tt.initialC
			vm.flags.Z = tt.initialZ

			// Execute the cmp function
			vm.cmp_B(nil)

			if vm.registers.A != tt.initialA {
				t.Error("Expected unchanged vm.registers.A")
			}

			if vm.registers.B != tt.initialB {
				t.Error("Expected unchanged vm.registers.B")
			}

			// Check flags
			if vm.flags.Z != tt.expectedZero {
				t.Errorf("%s: Expected Zero flag %t, got %t", tt.name, tt.expectedZero, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedSign {
				t.Errorf("%s: Expected Sign flag %t, got %t", tt.name, tt.expectedSign, vm.flags.S)
			}
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("%s: Expected Carry flag %t, got %t", tt.name, tt.expectedCarry, vm.flags.C)
			}
			if vm.flags.H != tt.expectedAuxCarry {
				t.Errorf("%s: Expected Aux Carry flag %t, got %t", tt.name, tt.expectedAuxCarry, vm.flags.H)
			}
			if vm.flags.P != tt.expectedParity {
				t.Errorf("%s: Expected Aux Parity flag %t, got %t", tt.name, tt.expectedParity, vm.flags.P)
			}
		})
	}
}
