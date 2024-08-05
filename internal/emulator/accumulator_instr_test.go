package emulator

import "testing"

func TestADC_B(t *testing.T) {
	tests := []struct {
		name         string
		initialA     byte
		initialB     byte
		initialCarry bool
		expectedA    byte
		expectedZ    bool
		expectedS    bool
		expectedC    bool
		expectedH    bool
		expectedP    bool
	}{
		{
			name:         "ADC without carry",
			initialA:     0x42,
			initialB:     0x3D,
			initialCarry: false,
			expectedA:    0x7F,
			expectedZ:    false,
			expectedS:    false,
			expectedC:    false,
			expectedH:    false,
			expectedP:    false,
		},
		{
			name:         "ADC with carry",
			initialA:     0x42,
			initialB:     0x3D,
			initialCarry: true,
			expectedA:    0x80,
			expectedZ:    false,
			expectedS:    true,
			expectedC:    false,
			expectedH:    true,
			expectedP:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := NewEmulator(&NullHardware{})
			vm.Registers.A = tt.initialA
			vm.Registers.B = tt.initialB
			vm.flags.C = tt.initialCarry

			vm.adc_B(nil)

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
			if vm.flags.H != tt.expectedH {
				t.Errorf("%s: expected H flag %t, got %t", tt.name, tt.expectedH, vm.flags.H)
			}
			if vm.flags.P != tt.expectedP {
				t.Errorf("%s: expected P flag %t, got %t", tt.name, tt.expectedP, vm.flags.P)
			}
		})
	}
}
func TestSub_A(t *testing.T) {

	vm := NewEmulator(&NullHardware{})
	vm.Registers.A = 0x3E
	vm.flags.C = true
	vm.flags.H = true
	vm.flags.S = true

	vm.sub_A(nil)

	if vm.Registers.A != 0x0 {
		t.Errorf("Expected A=0x0, got 0x%02X", vm.Registers.A)
	}
	if vm.flags.C {
		t.Errorf("Expected C flag false, got true")
	}
	if vm.flags.H {
		t.Errorf("Expected H flag false, got true")
	}
	if vm.flags.S {
		t.Errorf("Expected S flag false, got true")
	}
	if !vm.flags.Z {
		t.Errorf("Expected Z flag true, got false")
	}
	if !vm.flags.P {
		t.Errorf("Expected P flag true, got false")
	}
}

func TestSBB_L(t *testing.T) {
	vm := NewEmulator(&NullHardware{})
	vm.Registers.A = 4
	vm.Registers.L = 2
	vm.flags.C = true
	vm.flags.Z = true
	vm.flags.S = true
	vm.flags.P = true

	vm.sbb_L(nil)

	if vm.Registers.A != 1 {
		t.Errorf("Expected A=1, got 0x%02X", vm.Registers.A)
	}
	if vm.flags.C {
		t.Errorf("Expected C flag false, got true")
	}
	if !vm.flags.H {
		t.Errorf("Expected H flag true, got false")
	}
	if vm.flags.S {
		t.Errorf("Expected S flag false, got true")
	}
	if vm.flags.Z {
		t.Errorf("Expected Z flag true, got false")
	}
	if vm.flags.P {
		t.Errorf("Expected P flag false, got true")
	}
}

func TestXRA(t *testing.T) {
	tests := []struct {
		name      string
		initialA  byte
		reg       byte
		expectedA byte
		expectedZ bool
		expectedS bool
		expectedC bool
		expectedP bool
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
			expectedP: true,
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
			expectedParityFlag:  true,
		},
		{
			name:                "AND with itself",
			initialAccumulator:  0b10101010,
			data:                0b10101010,
			expectedAccumulator: 0b10101010,
			expectedZeroFlag:    false,
			expectedSignFlag:    true,
			expectedCarryFlag:   false,
			expectedParityFlag:  true,
		},
		{
			name:                "AND resulting in non-zero without sign",
			initialAccumulator:  0b11110000,
			data:                0b01111111,
			expectedAccumulator: 0b01110000,
			expectedZeroFlag:    false,
			expectedSignFlag:    false,
			expectedCarryFlag:   false,
			expectedParityFlag:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			vm := NewEmulator(&NullHardware{})
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
	vm := NewEmulator(&NullHardware{})

	tests := []struct {
		name      string
		initialA  byte
		data      byte
		expectedA byte
		expectedZ bool
		expectedS bool
		expectedC bool
		expectedP bool
	}{
		{
			name:      "ORA resulting in zero",
			initialA:  0x00,
			data:      0x00,
			expectedA: 0x00,
			expectedZ: true,
			expectedS: false,
			expectedC: false,
			expectedP: true,
		},
		{
			name:      "ORA resulting in non-zero",
			initialA:  0xF0,
			data:      0x0F,
			expectedA: 0xFF,
			expectedZ: false,
			expectedS: true,
			expectedC: false,
			expectedP: true,
		},
		{
			name:      "ORA with itself",
			initialA:  0x55,
			data:      0x55,
			expectedA: 0x55,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedP: true,
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

// TestCMP tests the cmp function which compares an 8-bit immediate value with the accumulator.
func TestCMP_B(t *testing.T) {
	// Initialize a CPU8080 instance with a dummy program and hardware IO
	vm := NewEmulator(&NullHardware{})

	// Test cases
	tests := []struct {
		name               string
		initialA           byte
		initialB           byte
		initialC, initialZ bool
		expectedZero       bool
		expectedSign       bool
		expectedCarry      bool
		expectedAuxCarry   bool
		expectedParity     bool
	}{
		{
			name:             "Accumulator greater",
			initialA:         0x0F,
			initialB:         0x07,
			expectedZero:     false,
			expectedSign:     false,
			expectedCarry:    false,
			expectedAuxCarry: false,
			expectedParity:   false,
		},
		{
			name:             "Accumulator equal",
			initialA:         0x12,
			initialB:         0x12,
			expectedZero:     true,
			expectedSign:     false,
			expectedCarry:    false,
			expectedAuxCarry: false,
			expectedParity:   true,
		},
		{
			name:             "docs 2",
			initialA:         0x02,
			initialB:         0x05,
			expectedZero:     false,
			expectedSign:     true,
			expectedCarry:    true,
			expectedAuxCarry: true,
			expectedParity:   false,
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
			vm.Registers.A = tt.initialA
			vm.Registers.B = tt.initialB
			vm.flags.C = tt.initialC
			vm.flags.Z = tt.initialZ

			// Execute the cmp function
			vm.cmp_B(nil)

			if vm.Registers.A != tt.initialA {
				t.Error("Expected unchanged vm.registers.A")
			}

			if vm.Registers.B != tt.initialB {
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
