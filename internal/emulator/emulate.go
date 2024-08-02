package emulator

import (
	"fmt"
	"time"
)

var stateCounts = []int{
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 00..0f
	4, 10, 7, 5, 5, 5, 7, 4, 4, 10, 7, 5, 5, 5, 7, 4, // 00..1f
	4, 10, 16, 5, 5, 5, 7, 4, 4, 10, 16, 5, 5, 5, 7, 4, // 20..2f
	4, 10, 13, 5, 10, 10, 10, 4, 4, 10, 13, 5, 5, 5, 7, 4, // 30..3f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 40..4f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 50..5f
	5, 5, 5, 5, 5, 5, 7, 5, 5, 5, 5, 5, 5, 5, 7, 5, // 60..6f
	7, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 7, 5, // 70..7f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 80..8f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // 90..9f
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // a0..af
	4, 4, 4, 4, 4, 4, 7, 4, 4, 4, 4, 4, 4, 4, 7, 4, // b0..bf
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 17, 7, 11, // c0..cf
	5, 10, 10, 10, 11, 11, 7, 11, 5, 10, 10, 10, 11, 17, 7, 11, // d0..df
	5, 10, 10, 18, 11, 11, 7, 11, 5, 5, 10, 5, 11, 17, 7, 11, // e0..ef
	5, 10, 10, 4, 11, 11, 7, 11, 5, 5, 10, 4, 11, 17, 7, 11, // f0..ff
}

func (vm *CPU8080) runCycles(cycleCount int) {
	var startTime time.Time
	if !vm.Options.UnlimitedTPS {
		startTime = time.Now()
	}

	for vm.cycleCount < cycleCount {
		select {
		case opcode := <-vm.interruptRequest:
			vm.handleInterrupt(opcode)
		default:
			if int(vm.pc) >= vm.programSize {
				break
			}
			currentCode := vm.memory[vm.pc : vm.pc+3]

			if vm.pc == 0x18DC {
				fmt.Printf("%v\n", vm.memory[0x2010:0x2020])
			}

			op := currentCode[0]
			vm.pc++
			vm.cycleCount += stateCounts[op]
			vm.totalCycles += stateCounts[op]

			if opcodeFunc, exists := vm.opcodeTable[op]; exists {
				opcodeFunc(currentCode[1:])
			} else {
				vm.Logger.Fatal("unsupported", "address", fmt.Sprintf("%04X", vm.pc-1), "opcode", fmt.Sprintf("%02X", op), "totalCycles", vm.totalCycles)
			}
		}
	}

	if !vm.Options.UnlimitedTPS {
		elapsed := time.Since(startTime)
		if remaining := (17 * time.Millisecond) - elapsed; remaining > 0 {
			time.Sleep(remaining)
		}
	}
}

func toUint16(high, low byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// NOP: No operation.
func (vm *CPU8080) nop(data []byte) {
	vm.Logger.Debugf("[00] NOP")
}

// OUT D8: Output accumulator to device at 8-bit immediate address.
func (vm *CPU8080) out(data []byte) {
	address := data[0]
	deviceName := vm.Hardware.OutDeviceName(address)
	vm.Logger.Debugf("[D3] OUT \t(%s),A", deviceName)
	vm.pc++
	err := vm.Hardware.Out(address, vm.registers.A)
	if err != nil {
		vm.Logger.Fatal("OUT", "address", fmt.Sprintf("%04X", vm.pc-2), "error", err)
	}
}

// IN D8: Input accumulator from device at 8-bit immediate address.
func (vm *CPU8080) in(data []byte) {
	address := data[0]
	deviceName := vm.Hardware.InDeviceName(address)
	vm.Logger.Debugf("[D8] IN  \tA,(%s)", deviceName)
	vm.pc++
	result, err := vm.Hardware.In(address)
	if err != nil {
		vm.Logger.Fatal("IN", "address", fmt.Sprintf("%04X", vm.pc-2), "error", err)
	}
	vm.registers.A = result
}
