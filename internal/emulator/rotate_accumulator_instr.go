package emulator

// RRC: Rotate accumulator right.
// The carry bit is set equal to the low-order
// bit of the accumulator. The contents of the accumulator are
// rotated one bit position to the right, with the low-order bit
// being transferred to the high-order bit position of the
// accumulator.
func (vm *CPU8080) rrc(data []byte) {
	vm.Logger.Debugf("[0F] RRC \tA")
	// Isolate least significant bit to check for Carry
	vm.flags.C = vm.Registers.A&0x01 == 1
	// Rotate accumulator right
	vm.Registers.A = (vm.Registers.A >> 1) | (vm.Registers.A << (8 - 1))
}

// RLC: Rotate accumulator left. The Carry bit is set equal to the high-order
// bit of the accumulator. The contents of the accumulator are rotated one bit
// position to the left, with the high-order bit being transferred to the
// low-order bit position of the accumulator
func (vm *CPU8080) rlc(data []byte) {
	vm.Logger.Debugf("[07] RLC \tA")
	// Isolate most significant bit to check for Carry
	vm.flags.C = (vm.Registers.A & 0x80) == 0x80
	// Rotate accumulator left
	vm.Registers.A = (vm.Registers.A << 1) | (vm.Registers.A >> (8 - 1))
}

// RAL: Rotate accumulator left through carry.
// The contents of the accumulator are rotated one bit position to the left.
// The high-order bit of the accumulator replaces the Carry bit, while the
// Carry bit replaces the high-order bit of the accumulator.
func (vm *CPU8080) ral(data []byte) {
	vm.Logger.Debugf("[17] RAL \tA")
	var carry uint8
	if vm.flags.C {
		carry = 1
	}
	// Isolate most significant bit to check for Carry
	vm.flags.C = (vm.Registers.A & 0x80) == 0x80
	// Rotate accumulator left through carry
	vm.Registers.A = (vm.Registers.A << 1) | carry
}

// RAR: Rotate accumulator right through carry.
// The contents of the accumulator are rotated one bit position to the right.
// The low order bit of the accumulator replaces the carry bit, while the carry bit replaces
// the high order bit of the accumulator.
func (vm *CPU8080) rar(data []byte) {
	vm.Logger.Debugf("[1F] RAR \tA")
	var carryRotate uint8
	if vm.flags.C {
		carryRotate = 1
	}
	// Isolate least significant bit to check for Carry
	vm.flags.C = vm.Registers.A&0x01 != 0
	// Rotate accumulator right through carry
	vm.Registers.A = (vm.Registers.A >> 1) | (carryRotate << (8 - 1))
}
