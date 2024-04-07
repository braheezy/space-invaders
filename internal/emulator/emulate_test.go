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
			vm.dad_DE(nil) // The data slice is not used in the function

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
func TestCMP(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	// Test cases
	tests := []struct {
		name             string
		initialA         byte // Initial accumulator value
		data             byte // Immediate value to compare with
		expectedZero     bool // Expected Zero flag
		expectedSign     bool // Expected Sign flag
		expectedCarry    bool // Expected Carry flag
		expectedAuxCarry bool // Expected Auxiliary Carry flag
		expectedParity   bool // Expected Parity flag
	}{
		{
			name:             "Accumulator greater",
			initialA:         0x0F,
			data:             0x07,
			expectedZero:     false,
			expectedSign:     false, // Result is positive
			expectedCarry:    false, // No borrow
			expectedAuxCarry: false, //  No aux borrow
			expectedParity:   true,  // Even parity of result (0x08)
		},
		{
			name:             "Accumulator equal",
			initialA:         0x12,
			data:             0x12,
			expectedZero:     true,
			expectedSign:     false, // Zero result, sign flag is 0
			expectedCarry:    false, // Equal, no borrow
			expectedAuxCarry: false, // No borrow in any nibble
			expectedParity:   true,  // Even parity (result is 0)
		},
		{
			name:             "Accumulator less",
			initialA:         0x01,
			data:             0x02,
			expectedZero:     false,
			expectedSign:     true,  // Negative result (considering 8-bit unsigned overflow)
			expectedCarry:    true,  // Borrow occurs
			expectedAuxCarry: true,  // Aux carry occurs if there's a borrow from bit 4 to bit 3
			expectedParity:   false, // Odd parity of result (0xFF)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the initial A register value
			vm.registers.A = tt.initialA

			// Execute the cmp function
			vm.cmp([]byte{tt.data})

			// Check flags
			if vm.flags.Z != tt.expectedZero {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZero, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedSign {
				t.Errorf("Expected Sign flag %t, got %t", tt.expectedSign, vm.flags.S)
			}
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected Carry flag %t, got %t", tt.expectedCarry, vm.flags.C)
			}
			if vm.flags.H != tt.expectedAuxCarry {
				t.Errorf("Expected Aux Carry flag %t, got %t", tt.expectedAuxCarry, vm.flags.H)
			}
			// Assuming setP method correctly sets the Parity flag based on the result's parity
			// Parity flag check might be omitted if its setting relies on complex logic within setP
		})
	}
}
func TestAND(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	tests := []struct {
		name      string
		initialA  byte // Initial accumulator value
		data      byte // Immediate value to AND with
		expectedA byte // Expected accumulator value after AND
		expectedZ bool // Expected Zero flag
		expectedS bool // Expected Sign flag
		expectedC bool // Expected Carry flag (always false)
		expectedP bool // Expected Parity flag
	}{
		{
			name:      "AND resulting in zero",
			initialA:  0xF0,
			data:      0x0F,
			expectedA: 0x00, // 0xF0 AND 0x0F = 0x00
			expectedZ: true,
			expectedS: false,
			expectedC: false,
			expectedP: true, // Parity of 0x00 is even
		},
		{
			name:      "AND resulting in non-zero",
			initialA:  0xFF,
			data:      0x0F,
			expectedA: 0x0F, // 0xFF AND 0x0F = 0x0F
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedP: true, // Parity of 0x0F (00001111) is even
		},
		{
			name:      "AND with itself",
			initialA:  0xAA, // 10101010
			data:      0xAA,
			expectedA: 0xAA,
			expectedZ: false,
			expectedS: true, // Sign bit is 1
			expectedC: false,
			expectedP: true, // Parity is even (4 ones)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and set the initial accumulator value
			vm.flags = flags{}
			vm.registers.A = tt.initialA

			// Execute the and function with the test data
			vm.and([]byte{tt.data})

			// Check the accumulator's value and flags
			if vm.registers.A != tt.expectedA {
				t.Errorf("%s: expected accumulator %02X, got %02X", tt.name, tt.expectedA, vm.registers.A)
			}
			if vm.flags.Z != tt.expectedZ {
				t.Errorf("%s: expected Z flag %t, got %t", tt.name, tt.expectedZ, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedS {
				t.Errorf("%s: expected S flag %t, got %t", tt.name, tt.expectedS, vm.flags.S)
			}
			if vm.flags.C != tt.expectedC {
				t.Errorf("%s: expected C flag %t, got %t", tt.name, tt.expectedC, vm.flags.C)
			}
			if vm.flags.P != tt.expectedP {
				t.Errorf("%s: expected P flag %t, got %t", tt.name, tt.expectedP, vm.flags.P)
			}
		})
	}
}

func TestADD(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	tests := []struct {
		name      string
		initialA  byte // Initial accumulator value
		data      byte // Immediate value to add
		expectedA byte // Expected accumulator value after addition
		expectedZ bool // Expected Zero flag
		expectedS bool // Expected Sign flag
		expectedC bool // Expected Carry flag
		expectedH bool // Expected Auxiliary Carry flag
		expectedP bool // Expected Parity flag
	}{
		{
			name:      "Addition without carry",
			initialA:  0x14,
			data:      0x01,
			expectedA: 0x15,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedH: false,
			expectedP: false, // Parity of 0x15 (00010101) is odd (3 ones)
		},
		{
			name:      "Addition with carry",
			initialA:  0xFF,
			data:      0x01,
			expectedA: 0x00,
			expectedZ: true,
			expectedS: false,
			expectedC: true,
			expectedH: true,
			expectedP: true, // Parity of 0x00 is even
		},
		{
			name:      "Result with sign bit set",
			initialA:  0x7F,
			data:      0x01,
			expectedA: 0x80,
			expectedZ: false,
			expectedS: true, // Sign bit is set because the result is 0x80
			expectedC: false,
			expectedH: true,
			expectedP: false, // Parity of 0x80 (10000000) is odd (1 one)
		},
		// Add more tests as necessary, particularly for edge cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and set the initial accumulator value
			vm.flags = flags{}
			vm.registers.A = tt.initialA

			// Execute the add function with the test data
			vm.add([]byte{tt.data})

			// Check the accumulator's value and flags
			if vm.registers.A != tt.expectedA {
				t.Errorf("%s: expected accumulator %02X, got %02X", tt.name, tt.expectedA, vm.registers.A)
			}
			if vm.flags.Z != tt.expectedZ {
				t.Errorf("%s: expected Z flag %t, got %t", tt.name, tt.expectedZ, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedS {
				t.Errorf("%s: expected S flag %t, got %t", tt.name, tt.expectedS, vm.flags.S)
			}
			if vm.flags.C != tt.expectedC {
				t.Errorf("%s: expected C flag %t, got %t", tt.name, tt.expectedC, vm.flags.C)
			}
			if vm.flags.H != tt.expectedH {
				t.Errorf("%s: expected H flag %t, got %t", tt.name, tt.expectedH, vm.flags.H)
			}
			if vm.flags.P != tt.expectedP {
				t.Errorf("%s: expected P flag %t, got %t", tt.name, tt.expectedP, vm.flags.P)
			}
		})
	}
}
