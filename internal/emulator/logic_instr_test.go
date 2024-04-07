package emulator

import "testing"

func TestXRA(t *testing.T) {
	tests := []struct {
		name      string
		initialA  byte // Initial value of the accumulator
		reg       byte // Value of the register to XOR with
		expectedA byte // Expected value of the accumulator after XOR
		expectedZ bool // Expected Zero flag
		expectedS bool // Expected Sign flag
		expectedC bool // Expected Carry flag (always false after XOR)
		expectedP bool // Expected Parity flag
	}{
		{
			name:      "Zero Result",
			initialA:  0xFF,
			reg:       0xFF,
			expectedA: 0x00,
			expectedZ: true,
			expectedS: false,
			expectedC: false,
			expectedP: true,
		},
		{
			name:      "Non-zero Result, Even Parity",
			initialA:  0xF0,
			reg:       0x0F,
			expectedA: 0xFF,
			expectedZ: false,
			expectedS: true,
			expectedC: false,
			expectedP: true, // 8 ones, even parity
		},
		{
			name:      "Non-zero Result, Odd Parity",
			initialA:  0b10101110,
			reg:       0b00110011,
			expectedA: 0b10011101,
			expectedZ: false,
			expectedS: true,
			expectedC: false,
			expectedP: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := &CPU8080{
				registers: struct{ A, B, C, D, E, H, L byte }{A: tt.initialA},
				flags:     flags{},
			}

			vm.xra(tt.reg)

			if vm.registers.A != tt.expectedA {
				t.Errorf("Expected accumulator value %02X, got %02X", tt.expectedA, vm.registers.A)
			}
			if vm.flags.Z != tt.expectedZ {
				t.Errorf("Expected Zero flag %t, got %t", tt.expectedZ, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedS {
				t.Errorf("Expected Sign flag %t, got %t", tt.expectedS, vm.flags.S)
			}
			if vm.flags.C != tt.expectedC {
				t.Errorf("Expected Carry flag %t, got %t", tt.expectedC, vm.flags.C)
			}
			if vm.flags.P != tt.expectedP {
				t.Errorf("Expected Parity flag %t, got %t", tt.expectedP, vm.flags.P)
			}
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
