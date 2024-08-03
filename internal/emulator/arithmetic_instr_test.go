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
			data:      0x42,
			expectedA: 0x56,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedH: false,
			expectedP: true,
		},
		{
			name:      "Addition with carry",
			initialA:  0x56,
			data:      0xBE,
			expectedA: 0x14,
			expectedZ: false,
			expectedS: false,
			expectedC: true,
			expectedH: true,
			expectedP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags and set the initial accumulator value
			vm.flags = flags{}
			vm.Registers.A = tt.initialA

			// Execute the add function with the test data
			vm.adi([]byte{tt.data})

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
			if vm.flags.H != tt.expectedH {
				t.Errorf("%s: expected H flag %t, got %t", tt.name, tt.expectedH, vm.flags.H)
			}
			if vm.flags.P != tt.expectedP {
				t.Errorf("%s: expected P flag %t, got %t", tt.name, tt.expectedP, vm.flags.P)
			}
		})
	}
}

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
			vm := NewCPU8080(&[]byte{}, nil)
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
			vm.Registers.H = tt.initialH
			vm.Registers.L = tt.initialL

			vm.dcx_H([]byte{0x00, 0x00})

			if vm.Registers.H != tt.expectedH || vm.Registers.L != tt.expectedL {
				t.Errorf("Expected H=0x%02X, L=0x%02X; got H=0x%02X, L=0x%02X", tt.expectedH, tt.expectedL, vm.Registers.H, vm.Registers.L)
			}
		})
	}
}

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
			vm := NewCPU8080(&[]byte{}, nil)
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

func TestIncC(t *testing.T) {
	vm := NewCPU8080(&[]byte{}, nil)
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

func TestSub_A(t *testing.T) {

	vm := NewCPU8080(&[]byte{}, nil)
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
	vm := NewCPU8080(&[]byte{}, nil)
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
