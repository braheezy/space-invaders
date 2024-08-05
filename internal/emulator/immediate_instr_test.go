package emulator

import "testing"

func TestADD(t *testing.T) {
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
		expectedH bool
		expectedP bool
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

func TestAND(t *testing.T) {
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
			name:      "AND resulting in zero",
			initialA:  0xF0,
			data:      0x0F,
			expectedA: 0x00,
			expectedZ: true,
			expectedS: false,
			expectedC: false,
			expectedP: true,
		},
		{
			name:      "AND resulting in non-zero",
			initialA:  0xFF,
			data:      0x0F,
			expectedA: 0x0F,
			expectedZ: false,
			expectedS: false,
			expectedC: false,
			expectedP: true,
		},
		{
			name:      "AND with itself",
			initialA:  0xAA,
			data:      0xAA,
			expectedA: 0xAA,
			expectedZ: false,
			expectedS: true,
			expectedC: false,
			expectedP: true,
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
