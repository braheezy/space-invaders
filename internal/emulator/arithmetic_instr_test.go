package emulator

import "testing"

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

func TestIncHL(t *testing.T) {
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
			initialH:      0x0F,
			initialL:      0xFE,
			expectedH:     0x0F,
			expectedL:     0xFF,
			carryFlagSet:  false,
			expectedCarry: false,
		},
		{
			name:          "Boundary increment with carry unchanged",
			initialH:      0xFF,
			initialL:      0xFF,
			expectedH:     0x00,
			expectedL:     0x00,
			carryFlagSet:  true,
			expectedCarry: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewCPU8080(&[]byte{}, nil) // Initialize your CPU8080 instance
			vm.registers.H = tt.initialH
			vm.registers.L = tt.initialL
			vm.flags.C = tt.carryFlagSet

			vm.inc_HL([]byte{0x01, 0x02}) // Execute the INC HL operation

			if vm.registers.H != tt.expectedH || vm.registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.registers.H, vm.registers.L)
			}
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("Expected Carry flag=%t; got %t", tt.expectedCarry, vm.flags.C)
			}
		})
	}
}
