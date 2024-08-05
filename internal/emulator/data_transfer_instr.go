package emulator

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

// MOV L,M: Load value from register B into memory address from register pair HL
func (vm *CPU8080) move_LM(data []byte) {
	vm.Logger.Debugf("[6E] LD  \tL,(HL)")
	vm.Registers.L = vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)]
}

// MOV D,B: Load value from register B into register D.
func (vm *CPU8080) move_DB(data []byte) {
	vm.Logger.Debugf("[50] LD  \tD,B")
	vm.Registers.D = vm.Registers.B
}

// MOV D,E: Load value from register E into register D.
func (vm *CPU8080) move_DE(data []byte) {
	vm.Logger.Debugf("[53] LD  \tD,E")
	vm.Registers.D = vm.Registers.E
}

// MOV E,B: Load value from register B into register E.
func (vm *CPU8080) move_EB(data []byte) {
	vm.Logger.Debugf("[58] LD  \tE,B")
	vm.Registers.E = vm.Registers.B
}

// MOV E,L: Load value from register L into register E.
func (vm *CPU8080) move_EL(data []byte) {
	vm.Logger.Debugf("[5D] LD  \tE,L")
	vm.Registers.E = vm.Registers.L
}

// MOV B,A: Load value from accumulator into register B.
func (vm *CPU8080) move_BA(data []byte) {
	vm.Logger.Debugf("[47] LD  \tB,A")
	vm.Registers.B = vm.Registers.A
}

// MOV B,D: Load value from register B into register D.
func (vm *CPU8080) move_BD(data []byte) {
	vm.Logger.Debugf("[42] LD  \tB,D")
	vm.Registers.B = vm.Registers.D
}

// MOV B,E: Load value from register B into register E.
func (vm *CPU8080) move_BE(data []byte) {
	vm.Logger.Debugf("[43] LD  \tB,E")
	vm.Registers.B = vm.Registers.E
}

// MOV C,A: Load value from accumulator into register C.
func (vm *CPU8080) move_CA(data []byte) {
	vm.Logger.Debugf("[4F] LD  \tC,A")
	vm.Registers.C = vm.Registers.A
}

// MOV C,B: Load value from register B into register C.
func (vm *CPU8080) move_CB(data []byte) {
	vm.Logger.Debugf("[48] LD  \tC,B")
	vm.Registers.C = vm.Registers.B
}

// MOV C,D: Load value from register D into register C.
func (vm *CPU8080) move_CD(data []byte) {
	vm.Logger.Debugf("[4A] LD  \tC,D")
	vm.Registers.C = vm.Registers.D
}

// MOV C,E: Load value from register E into register C.
func (vm *CPU8080) move_CE(data []byte) {
	vm.Logger.Debugf("[4B] LD  \tC,E")
	vm.Registers.C = vm.Registers.E
}

// MOV C,H: Load value from register H into register C.
func (vm *CPU8080) move_CH(data []byte) {
	vm.Logger.Debugf("[4C] LD  \tC,H")
	vm.Registers.C = vm.Registers.H
}

// MOV H,B: Load value from register B into register H.
func (vm *CPU8080) move_HB(data []byte) {
	vm.Logger.Debugf("[60] LD  \tH,B")
	vm.Registers.H = vm.Registers.B
}

// MOV H,L: Load value from register L into register H.
func (vm *CPU8080) move_HL(data []byte) {
	vm.Logger.Debugf("[65] LD  \tH,L")
	vm.Registers.H = vm.Registers.L
}

// MOV A,C: Load value from register C into accumulator.
func (vm *CPU8080) move_AC(data []byte) {
	vm.Logger.Debugf("[79] LD  \tA,C")
	vm.Registers.A = vm.Registers.C
}

// MOV D,C: Load value from register C into register D.
func (vm *CPU8080) move_DC(data []byte) {
	vm.Logger.Debugf("[51] LD  \tD,C")
	vm.Registers.D = vm.Registers.C
}

// MOV D,H: Load value from register H into register D.
func (vm *CPU8080) move_DH(data []byte) {
	vm.Logger.Debugf("[54] LD  \tD,H")
	vm.Registers.D = vm.Registers.H
}

// MOV D,L: Load value from register L into register D.
func (vm *CPU8080) move_DL(data []byte) {
	vm.Logger.Debugf("[55] LD  \tD,L")
	vm.Registers.D = vm.Registers.L
}

// MOV H,C: Load value from register C into register H.
func (vm *CPU8080) move_HC(data []byte) {
	vm.Logger.Debugf("[61] LD  \tH,C")
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
func (vm *CPU8080) move_MB(data []byte) {
	vm.Logger.Debugf("[70] LD  \t(HL),B")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.B
}

// MOV M,C: Move register C into memory location pointed to by register pair HL.
func (vm *CPU8080) move_MC(data []byte) {
	vm.Logger.Debugf("[71] LD  \t(HL),C")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.C
}

// MOV M,D: Move register D into memory location pointed to by register pair HL.
func (vm *CPU8080) move_MD(data []byte) {
	vm.Logger.Debugf("[72] LD  \t(HL),D")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.D
}

// MOV M,E: Move register E into memory location pointed to by register pair HL.
func (vm *CPU8080) move_ME(data []byte) {
	vm.Logger.Debugf("[73] LD  \t(HL),E")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.E
}

// MOV M,H: Move register H into memory location pointed to by register pair HL.
func (vm *CPU8080) move_MH(data []byte) {
	vm.Logger.Debugf("[74] LD  \t(HL),H")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.H
}

// MOV M,L: Move register L into memory location pointed to by register pair HL.
func (vm *CPU8080) move_ML(data []byte) {
	vm.Logger.Debugf("[75] LD  \t(HL),L")
	vm.Memory[toUint16(vm.Registers.H, vm.Registers.L)] = vm.Registers.L
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

// MOV B,C: Move value from register C into register B.
func (vm *CPU8080) move_BC(data []byte) {
	vm.Logger.Debugf("[41] LD  \tB,C")
	vm.Registers.B = vm.Registers.C
}

// MOV B,L: Move value from register L into register B.
func (vm *CPU8080) move_BL(data []byte) {
	vm.Logger.Debugf("[45] LD  \tB,L")
	vm.Registers.B = vm.Registers.L
}

// MOV B,H: Move value from register H into register B.
func (vm *CPU8080) move_BH(data []byte) {
	vm.Logger.Debugf("[44] LD  \tB,H")
	vm.Registers.B = vm.Registers.H
}

// MOV C,L: Move value from register C into register L.
func (vm *CPU8080) move_CL(data []byte) {
	vm.Logger.Debugf("[4D] LD  \tC,L")
	vm.Registers.C = vm.Registers.L
}

// MOV A,D: Move value from register D into accumulator.
func (vm *CPU8080) move_AD(data []byte) {
	vm.Logger.Debugf("[7A] LD  \tA,D")
	vm.Registers.A = vm.Registers.D
}

// MOV E,D: Move value from register D into register E.
func (vm *CPU8080) move_ED(data []byte) {
	vm.Logger.Debugf("[5A] LD  \tE,D")
	vm.Registers.E = vm.Registers.D
}

// MOV E,H: Move value from register H into register E.
func (vm *CPU8080) move_EH(data []byte) {
	vm.Logger.Debugf("[5C] LD  \tE,H")
	vm.Registers.E = vm.Registers.H
}

// MOV H,D: Move value from register D into register H.
func (vm *CPU8080) move_HD(data []byte) {
	vm.Logger.Debugf("[62] LD  \tH,D")
	vm.Registers.H = vm.Registers.D
}

// MOV L,C: Move value from register C into register L.
func (vm *CPU8080) move_LC(data []byte) {
	vm.Logger.Debugf("[69] LD  \tL,C")
	vm.Registers.L = vm.Registers.C
}

// MOV L,D: Move value from register D into register L.
func (vm *CPU8080) move_LD(data []byte) {
	vm.Logger.Debugf("[6A] LD  \tL,D")
	vm.Registers.L = vm.Registers.D
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

// MOV H,E: Move value from register E into register H.
func (vm *CPU8080) move_HE(data []byte) {
	vm.Logger.Debugf("[63] LD  \tH,E")
	vm.Registers.H = vm.Registers.E
}

// MOV E,C: Move value from register C into register E.
func (vm *CPU8080) move_EC(data []byte) {
	vm.Logger.Debugf("[59] LD  \tE,C")
	vm.Registers.E = vm.Registers.C
}

// MOV L,E: Move value from register E into register L.
func (vm *CPU8080) move_LE(data []byte) {
	vm.Logger.Debugf("[6B] LD  \tL,E")
	vm.Registers.L = vm.Registers.E
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

// MOV L,H: Move value from register H into register L.
func (vm *CPU8080) move_LH(data []byte) {
	vm.Logger.Debug("[6C] LD  \tH,L")
	vm.Registers.L = vm.Registers.H
}

// MOV D,A: Move value from accumulator into register D.
func (vm *CPU8080) move_DA(data []byte) {
	vm.Logger.Debug("[57] LD  \tD,A")
	vm.Registers.D = vm.Registers.A
}

// STAX B: Store accumulator in 16-bit immediate address pointed to by register pair BC
func (vm *CPU8080) stax_B(data []byte) {
	address := toUint16(vm.Registers.B, vm.Registers.C)
	vm.Logger.Debug("[32] LD  \t(BC),A")
	vm.Memory[address] = vm.Registers.A
}

// STAX D: Store accumulator in 16-bit immediate address pointed to by register pair DE
func (vm *CPU8080) stax_D(data []byte) {
	address := toUint16(vm.Registers.D, vm.Registers.E)
	vm.Logger.Debug("[12] LD  \t(DE),A")
	vm.Memory[address] = vm.Registers.A
}
