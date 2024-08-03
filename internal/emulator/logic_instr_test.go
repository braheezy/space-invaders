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
				Registers: struct{ A, B, C, D, E, H, L byte }{A: tt.initialA},
				flags:     flags{},
			}

			vm.xra(tt.reg)

			if vm.Registers.A != tt.expectedA {
				t.Errorf("Expected accumulator value %02X, got %02X", tt.expectedA, vm.Registers.A)
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
			vm.Registers.A = tt.initialA

			// Execute the and function with the test data
			vm.and([]byte{tt.data})

			// Check the accumulator's value and flags
			if vm.Registers.A != tt.expectedA {
				t.Errorf("%s: expected accumulator %02X, got %02X", tt.name, tt.expectedA, vm.Registers.A)
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

func TestANA(t *testing.T) {
	tests := []struct {
		name                string
		initialAccumulator  byte
		data                byte
		expectedAccumulator byte
		expectedZeroFlag    bool
		expectedSignFlag    bool
		expectedCarryFlag   bool
		expectedParityFlag  bool
	}{
		{
			name:                "AND with zero",
			initialAccumulator:  0b10101010,
			data:                0b00000000,
			expectedAccumulator: 0b00000000,
			expectedZeroFlag:    true,
			expectedSignFlag:    false,
			expectedCarryFlag:   false,
			expectedParityFlag:  true, // Parity is even because the result has an even number of 1 bits (zero in this case).
		},
		{
			name:                "AND with itself",
			initialAccumulator:  0b10101010,
			data:                0b10101010,
			expectedAccumulator: 0b10101010,
			expectedZeroFlag:    false,
			expectedSignFlag:    true, // Sign bit is 1 because the result's most significant bit is 1.
			expectedCarryFlag:   false,
			expectedParityFlag:  true, // Parity is even  because the result has an even number of 1 bits.
		},
		{
			name:                "AND resulting in non-zero without sign",
			initialAccumulator:  0b11110000,
			data:                0b01111111,
			expectedAccumulator: 0b01110000,
			expectedZeroFlag:    false,
			expectedSignFlag:    false,
			expectedCarryFlag:   false,
			expectedParityFlag:  false, // Parity is odd.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			vm := NewCPU8080(&[]byte{}, nil) // Assuming NewCPU8080 initializes the CPU state including flags and registers
			vm.Registers.A = tt.initialAccumulator

			// Execute
			vm.ana(tt.data)

			// Verify accumulator and flags
			if vm.Registers.A != tt.expectedAccumulator {
				t.Errorf("Expected accumulator %02x, got %02x", tt.expectedAccumulator, vm.Registers.A)
			}
			if vm.flags.Z != tt.expectedZeroFlag {
				t.Errorf("Expected zero flag %t, got %t", tt.expectedZeroFlag, vm.flags.Z)
			}
			if vm.flags.S != tt.expectedSignFlag {
				t.Errorf("Expected sign flag %t, got %t", tt.expectedSignFlag, vm.flags.S)
			}
			if vm.flags.C != tt.expectedCarryFlag {
				t.Errorf("Expected carry flag %t, got %t", tt.expectedCarryFlag, vm.flags.C)
			}
			if vm.flags.P != tt.expectedParityFlag {
				t.Errorf("Expected parity flag %t, got %t", tt.expectedParityFlag, vm.flags.P)
			}
		})
	}
}

func TestORA(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	tests := []struct {
		name      string
		initialA  byte // Initial accumulator value
		data      byte // Immediate value to OR with
		expectedA byte // Expected accumulator value after OR
		expectedZ bool // Expected Zero flag
		expectedS bool // Expected Sign flag
		expectedC bool // Expected Carry flag (always false)
		expectedP bool // Expected Parity flag
	}{
		{
			name:      "ORA resulting in zero",
			initialA:  0x00,
			data:      0x00,
			expectedA: 0x00, // 0x00 OR 0x00 = 0x00
			expectedZ: true,
			expectedS: false,
			expectedC: false,
			expectedP: true, // Parity of 0x00 is even
		},
		{
			name:      "ORA resulting in non-zero",
			initialA:  0xF0,
			data:      0x0F,
			expectedA: 0xFF, // 0xF0 OR 0x0F = 0xFF
			expectedZ: false,
			expectedS: true, // Sign bit is 1
			expectedC: false,
			expectedP: true, // Parity of 0xFF (11111111) is even
		},
		{
			name:      "ORA with itself",
			initialA:  0x55, // 01010101
			data:      0x55,
			expectedA: 0x55,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedP: true, // Parity is even (4 ones)
		},
		{
			name:      "ORA C example from docs",
			initialA:  0x33,
			data:      0x0F,
			expectedA: 0x3F,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and set the initial accumulator value
			vm.flags = flags{}
			vm.Registers.A = tt.initialA

			// Execute the ora function with the test data
			vm.ora(tt.data)

			// Check the accumulator's value and flags
			if vm.Registers.A != tt.expectedA {
				t.Errorf("%s: expected accumulator %02X, got %02X", tt.name, tt.expectedA, vm.Registers.A)
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

func TestRLC(t *testing.T) {

	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewCPU8080(&[]byte{}, nil)

	tests := []struct {
		name      string
		initialA  byte // Initial accumulator value
		expectedA byte // Expected accumulator value after RLC
		expectedC bool // Expected Carry flag
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
	vm := NewCPU8080(&[]byte{}, nil)
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

func TestXCHG(t *testing.T) {

	vm := NewCPU8080(&[]byte{}, nil)
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

func TestRAL(t *testing.T) {

	vm := NewCPU8080(&[]byte{}, nil)
	vm.Registers.A = 0xB5

	vm.ral(nil)

	if vm.Registers.A != 0x6A {
		t.Errorf("Expected accumulator to be 0x6A, got %02X", vm.Registers.A)
	}
	if !vm.flags.C {
		t.Error("Expected carry flag to be set")
	}
}
