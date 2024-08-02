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
			vm.adi([]byte{tt.data})

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
			vm := NewCPU8080(&[]byte{}, nil)
			vm.registers.H = tt.initialH
			vm.registers.L = tt.initialL
			vm.flags.C = tt.carryFlagSet

			vm.inx_H([]byte{0x00, 0x00})

			if vm.registers.H != tt.expectedH || vm.registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.registers.H, vm.registers.L)
			}
		})
	}
}

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
			vm := NewCPU8080(&[]byte{}, nil)
			vm.registers.H = tt.initialH
			vm.registers.L = tt.initialL
			vm.memory[tt.memoryLocation] = tt.initialMemoryValue

			vm.dcr_M([]byte{0x00, 0x00})

			if vm.memory[tt.memoryLocation] != tt.expectedMemoryValue {
				t.Errorf("Expected memory value=0x%02X; got 0x%02X", tt.expectedMemoryValue, vm.memory[tt.memoryLocation])
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
			vm := NewCPU8080(&[]byte{}, nil)
			vm.registers.H = tt.initialH
			vm.registers.L = tt.initialL
			vm.registers.D = tt.initialD
			vm.registers.E = tt.initialE
			vm.flags.C = true

			vm.dad_D([]byte{0x00, 0x00})

			if vm.registers.H != tt.expectedH || vm.registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.registers.H, vm.registers.L)
			}
			if vm.flags.C != tt.expectedCarry {
				t.Errorf("expected C flag %t, got %t", tt.expectedCarry, vm.flags.C)
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
			vm := NewCPU8080(&[]byte{}, nil)
			vm.registers.H = tt.initialH
			vm.registers.L = tt.initialL

			vm.dcx_H([]byte{0x00, 0x00})

			if vm.registers.H != tt.expectedH || vm.registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.registers.H, vm.registers.L)
			}
		})
	}
}
