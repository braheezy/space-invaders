package emulator

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type HardwareIO interface {
	In(addr byte) (byte, error)
	Out(addr byte, data byte) error
	InDeviceName(port byte) string
	OutDeviceName(port byte) string
	InterruptConditions() []InterruptCondition
	CyclesPerFrame() int
	Draw(*ebiten.Image)
	Init(*[65536]byte)
	Width() int
	Height() int
	Scale() int
	StartAddress() int
}
