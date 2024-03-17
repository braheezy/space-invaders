package emulator

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
)

type CPU8080 struct {
	// pc is the current program counter, the address of the next instruction to be executed
	pc uint16
	// programData is a pointer to bytes containing a sequence of 8080 instructions
	programData *[]byte
	// programSize is the number of bytes in the program
	programSize int
	// registers are the CPU's 16-bit registers
	registers registers
	// Logger object to use
	Logger *log.Logger
}

type registers struct {
	A, B, C, D, E, H, L byte
}

func NewCPU8080(program *[]byte) *CPU8080 {
	vm := &CPU8080{
		programData: program,
		programSize: len(*program),
		Logger:      log.New(os.Stdout),
	}

	return vm
}

func (vm *CPU8080) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	vm.nextOpCode()
	if int(vm.pc) == vm.programSize {
		return ebiten.Termination
	}

	return nil
}

func (vm *CPU8080) Draw(screen *ebiten.Image) {

}
func (vm *CPU8080) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
