package emulator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"

	"github.com/charmbracelet/log"
	"github.com/ebitengine/oto/v3"
)

type SoundManager struct {
	ctx     *oto.Context
	players map[string]*oto.Player
}

func NewSoundManager(sampleRate int, channelCount int, soundFiles embed.FS) (*SoundManager, error) {
	ctx, ready, err := oto.NewContext(
		&oto.NewContextOptions{
			// Typically 44100 or 48000
			SampleRate: sampleRate,
			// only 1 or 2 are supported by oto
			ChannelCount: channelCount,
			Format:       oto.FormatSignedInt16LE,
		})
	if err != nil {
		return nil, fmt.Errorf("oto.NewContext failed: " + err.Error())
	}

	sm := &SoundManager{}
	// Initialize sound players, one per unique sound
	sm.players = make(map[string]*oto.Player)
	sm.ctx = ctx

	// Wait for the audio context to be ready
	<-ready

	err = fs.WalkDir(soundFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			data, err := soundFiles.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			player := sm.ctx.NewPlayer(bytes.NewReader(data))
			sm.players[path] = player
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func NewSoundManagerWithDefaults(soundFiles embed.FS) (*SoundManager, error) {
	return NewSoundManager(44100, 2, soundFiles)
}

func (sm *SoundManager) Play(filePath string) {
	if player, exists := sm.players[filePath]; exists {
		player.Play()
	}
}

func (sm *SoundManager) Pause(filePath string) {
	if player, exists := sm.players[filePath]; exists && player.IsPlaying() {
		player.Pause()
	}
}
