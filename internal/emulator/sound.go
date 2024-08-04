package emulator

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/braheezy/qoa"
	"github.com/charmbracelet/log"
	"github.com/ebitengine/oto/v3"
	"github.com/go-audio/wav"
)

type SoundManager struct {
	ctx     *oto.Context
	players map[string]*oto.Player
}

func NewSoundManager(sampleRate int, channelCount int, soundFiles embed.FS) (*SoundManager, error) {
	ctx, ready, err := oto.NewContext(
		&oto.NewContextOptions{
			// Typically 44100 or 48000
			SampleRate: 44100,
			// only 1 or 2 are supported by oto
			ChannelCount: 1,
			Format:       oto.FormatSignedInt16LE,
		})
	if err != nil {
		return nil, fmt.Errorf("oto.NewContext failed: " + err.Error())
	}

	// Wait for the audio context to be ready
	<-ready

	sm := &SoundManager{}
	// Initialize sound players, one per unique sound
	sm.players = make(map[string]*oto.Player)
	sm.ctx = ctx

	err = fs.WalkDir(soundFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			data, err := soundFiles.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			var player *oto.Player
			switch filepath.Ext(path) {
			case ".wav":
				player, err = setupWavPlayer(data, sm.ctx)
				if err != nil {
					log.Fatal(err)
				}
			case ".qoa":
				player, err = setupQoaPlayer(data, sm.ctx)
				if err != nil {
					log.Fatal(err)
				}
			}
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
		player.Seek(0, io.SeekStart)
		player.Play()
	}
}

func (sm *SoundManager) Pause(filePath string) {
	if player, exists := sm.players[filePath]; exists && player.IsPlaying() {
		player.Pause()
	}
}

func setupWavPlayer(data []byte, ctx *oto.Context) (*oto.Player, error) {
	wavReader := bytes.NewReader(data)
	wavDecoder := wav.NewDecoder(wavReader)
	if !wavDecoder.IsValidFile() {
		return nil, errors.New("invalid WAV file")
	}

	pcmBuffer, err := wavDecoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	// Convert to bytes.
	byteData := make([]byte, len(pcmBuffer.Data))
	for i, sample := range pcmBuffer.Data {
		byteData[i] = byte(sample & 0xFF)
	}

	return ctx.NewPlayer(bytes.NewReader(byteData)), nil
}

func setupQoaPlayer(data []byte, ctx *oto.Context) (*oto.Player, error) {
	qoaMetadata, qoaAudioData, err := qoa.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding QOA data: %v", err)
	}

	reader := qoa.NewReader(qoaAudioData, int(qoaMetadata.Channels))
	return ctx.NewPlayer(reader), nil
}
