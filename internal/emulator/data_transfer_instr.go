package emulator

// LXI SP, D16: Load 16-bit immediate value into register pair SP.
func (vm *CPU8080) load_SP(data []byte) {
	operand := toUint16(data[1], data[0])
	vm.Logger.Debugf("[31] LD  \tSP,$%04X", operand)
	vm.sp = operand
	vm.PC += 2
}

// LXI B, D16: Load 16-bit immediate value into register pair B.
func (vm *CPU8080) load_BC(data []byte) {
	vm.Logger.Debugf("[01] LD  \tB,$%04X", toUint16(data[1], data[0]))
	vm.Registers.C = data[0]
	vm.Registers.B = data[1]
	vm.PC += 2
}

// LXI D, D16: Load 16-bit immediate value into register pair D.
func (vm *CPU8080) load_DE(data []byte) {
	vm.Logger.Debugf("[11] LD  \tDE,$%04X", toUint16(data[1], data[0]))
	vm.Registers.E = data[0]
	vm.Registers.D = data[1]
	vm.PC += 2
}

// LXI H, D16: Load 16-bit immediate value into register pair H.
func (vm *CPU8080) load_HL(data []byte) {
	vm.Logger.Debugf("[21] LD  \tHL,$%04X", toUint16(data[1], data[0]))
	vm.Registers.L = data[0]
	vm.Registers.H = data[1]
	vm.PC += 2
}

// MVI A, D8: Move 8-bit immediate value into accumulator.
func (vm *CPU8080) moveImm_A(data []byte) {
	vm.Logger.Debugf("[3E] LD  \tA,$%02X", data[0])
	vm.Registers.A = data[0]
	vm.PC++
}

// MVI B, D8: Move 8-bit immediate value into register B.
func (vm *CPU8080) moveImm_B(data []byte) {
	vm.Logger.Debugf("[06] LD  \tB,$%02X", data[0])
	vm.Registers.B = data[0]
	vm.PC++
}

// MVI C, D8: Move 8-bit immediate value into register C.
func (vm *CPU8080) moveImm_C(data []byte) {
	vm.Logger.Debugf("[0E] LD  \tC,$%02X", data[0])
	vm.Registers.C = data[0]
	vm.PC++
}

// MVI H, D8: Move 8-bit immediate value into register H.
func (vm *CPU8080) moveImm_H(data []byte) {
	vm.Logger.Debugf("[26] LD  \tH,$%02X", data[0])
	vm.Registers.H = data[0]
	vm.PC++
}

// MVI L, D8: Move 8-bit immediate value into register L.
func (vm *CPU8080) moveImm_L(data []byte) {
	vm.Logger.Debugf("[2E] LD  \tL,$%02X", data[0])
	vm.Registers.L = data[0]
	vm.PC++
}

// MVI D, D8: Move 8-bit immediate value into register L.
func (vm *CPU8080) moveImm_D(data []byte) {
	vm.Logger.Debugf("[16] LD  \tD,$%02X", data[0])
	vm.Registers.D = data[0]
	vm.PC++
}

// MVI M: Move 8-bit immediate value into memory address from register pair HL
func (vm *CPU8080) moveImm_M(data []byte) {
	address := toUint16(vm.Registers.H, vm.Registers.L)
	vm.Logger.Debugf("[36] LD  \t(HL),$%02X", data[0])
	vm.Memory[address] = data[0]
	vm.PC++
}

// LDAX D: Load value from address in register pair D into accumulator.
func (vm *CPU8080) loadAddr_D(data []byte) {
	vm.Logger.Debugf("[1A] LD  \tA,(DE)")
	vm.Registers.A = vm.Memory[toUint16(vm.Registers.D, vm.Registers.E)]
}

// LDAX B: Load value from address in register pair B into accumulator.
func (vm *CPU8080) loadAddr_B(data []byte) {
	vm.Logger.Debugf("[0A] LD  \tA,(BC)")
	vm.Registers.A = vm.Memory[toUint16(vm.Registers.B, vm.Registers.C)]
}

// MOV M,A: Move value from accumulator into register pair H.
func (vm *CPU8080) move_MA(data []byte) {
	address := toUint16(vm.Registers.H, vm.Registers.L)
	vm.Logger.Debugf("[77] LD  \t(HL),A ($%04X)", address)
	vm.Memory[address] = vm.Registers.A
}

// MOV L,A: Load value from accumulator into register L.
func (vm *CPU8080) move_LA(data []byte) {
	vm.Logger.Debugf("[6F] LD  \tL,A")
	vm.Registers.L = vm.Registers.A
}

// MOV L,B: Load value from register B into register L.
func (vm *CPU8080) move_LB(data []byte) {
	vm.Logger.Debugf("[68] LD  \tL,B")
	vm.Registers.L = vm.Registers.B
}

// MOV B,A: Load value from register A into register B.
func (vm *CPU8080) move_BA(data []byte) {
	vm.Logger.Debugf("[47] LD  \tB,A")
	vm.Registers.B = vm.Registers.A
}

// MOV C,A: Load value from accumulator into register C.
func (vm *CPU8080) move_CA(data []byte) {
	vm.Logger.Debugf("[4F] LD  \tC,A")
	vm.Registers.C = vm.Registers.A
}

// MOV A,C: Load value from register C into accumulator.
func (vm *CPU8080) move_AC(data []byte) {
	vm.Logger.Debugf("[79] LD  \tA,C")
	vm.Registers.A = vm.Registers.C
}

// MOV H,C: Load value from register C into register H.
func (vm *CPU8080) move_HC(data []byte) {
	vm.Logger.Debugf("[61] LD  \tA,C")
	vm.Registers.H = vm.Registers.C
}

// MOV E,M: Move memory location pointed to by register pair HL into register E.
func (vm *CPU8080) move_EM(data []byte) {
	vm.Logger.Debugf("[5E] LD  \tE,(HL)")
	vm.Registers.E = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV B,M: Move memory location pointed to by register pair HL into register B.
func (vm *CPU8080) move_BM(data []byte) {
	vm.Logger.Debugf("[46] LD  \tB,(HL)")
	vm.Registers.B = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV C,M: Move memory location pointed to by register pair HL into register C.
func (vm *CPU8080) move_CM(data []byte) {
	vm.Logger.Debugf("[4E] LD  \tC,(HL)")
	vm.Registers.C = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV D,M: Move memory location pointed to by register pair HL into register D.
func (vm *CPU8080) move_DM(data []byte) {
	vm.Logger.Debugf("[56] LD  \tD,(HL)")
	vm.Registers.D = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV A,M: Move memory location pointed to by register pair HL into register A.
func (vm *CPU8080) move_AM(data []byte) {
	vm.Logger.Debugf("[7E] LD  \tA,(HL)")
	vm.Registers.A = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV H,M: Move memory location pointed to by register pair HL into register H.
func (vm *CPU8080) move_HM(data []byte) {
	vm.Logger.Debugf("[66] LD  \tH,(HL)")
	vm.Registers.H = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV M,B: Move register B into memory location pointed to by register pair HL.
func (vm *CPU8080) move_MH(data []byte) {
	vm.Logger.Debugf("[70] LD  \t(HL),B")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.B
}

// MOV A,H: Move value from register H into accumulator.
func (vm *CPU8080) move_AH(data []byte) {
	vm.Logger.Debugf("[7C] LD  \tA,H")
	vm.Registers.A = vm.Registers.H
}

// MOV A,L: Move value from register L into accumulator.
func (vm *CPU8080) move_AL(data []byte) {
	vm.Logger.Debugf("[7D] LD  \tA,L")
	vm.Registers.A = vm.Registers.L
}

// MOV A,D: Move value from register D into accumulator.
func (vm *CPU8080) move_AD(data []byte) {
	vm.Logger.Debugf("[7A] LD  \tA,D")
	vm.Registers.A = vm.Registers.D
}

// MOV A,E: Move value from register E into accumulator.
func (vm *CPU8080) move_AE(data []byte) {
	vm.Logger.Debugf("[7B] LD  \tA,E")
	vm.Registers.A = vm.Registers.E
}

// MOV H,A: Move value from accumulator into register H.
func (vm *CPU8080) move_HA(data []byte) {
	vm.Logger.Debugf("[67] LD  \tH,A")
	vm.Registers.H = vm.Registers.A
}

// MOV A,B: Move value from register B into accumulator.
func (vm *CPU8080) move_AB(data []byte) {
	vm.Logger.Debugf("[78] LD  \tA,B")
	vm.Registers.A = vm.Registers.B
}

// MOV E,A: Move value from accumulator into register E.
func (vm *CPU8080) move_EA(data []byte) {
	vm.Logger.Debug("[5F] LD  \tE,A")
	vm.Registers.E = vm.Registers.A
}

// MOV D,A: Move value from accumulator into register D.
func (vm *CPU8080) move_DA(data []byte) {
	vm.Logger.Debug("[57] LD  \tD,A")
	vm.Registers.D = vm.Registers.A
}

// STA A16: Store accumulator in 16-bit immediate address.
func (vm *CPU8080) store_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[32] LD  \t$%04X,A", address)
	vm.Memory[address] = vm.Registers.A
	vm.PC += 2
}

// STAX B: Store accumulator in 16-bit immediate address.
func (vm *CPU8080) stax_B(data []byte) {
	address := toUint16(vm.Registers.B, vm.Registers.C)
	vm.Logger.Debug("[32] LD  \t(BC),A")
	vm.Memory[address] = vm.Registers.A
	vm.PC += 2
}

// LDA A16: Load accumulator from 16-bit immediate address.
func (vm *CPU8080) load_A(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[3A] LD  \tA,$%04X", address)
	vm.Registers.A = vm.Memory[address]
	vm.PC += 2
}

// LHLD A16: Load register pair HL from 16-bit immediate address.
func (vm *CPU8080) loadImm_HL(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[2A] LD  \tHL,$%04X", address)
	vm.Registers.L = vm.Memory[address]
	vm.Registers.H = vm.Memory[address+1]
	vm.PC += 2
}

// SHLD A16: Store register pair HL into 16-bit immediate address.
func (vm *CPU8080) store_HL(data []byte) {
	address := toUint16(data[1], data[0])
	vm.Logger.Debugf("[22] LD  \t$%04X,HL", address)
	vm.Memory[address] = vm.Registers.L
	vm.Memory[address+1] = vm.Registers.H
	vm.PC += 2
}
