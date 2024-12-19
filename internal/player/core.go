package player

import (
	"fmt"
	"github.com/ebitengine/oto/v3"
	"github.com/happystraw/text-player/internal/config"
	"io"
	"sync"
	"time"
)

type Player interface {
	Play(r io.ReadCloser) error
}

type player struct {
	ctx         *oto.Context
	initialized bool
	sync.Mutex
}

var p = &player{}

func (p *player) init(cfg *config.Engine) error {
	if p.initialized {
		return nil
	}

	p.Lock()
	defer p.Unlock()
	op := &oto.NewContextOptions{
		SampleRate:   cfg.SampleRate,
		ChannelCount: cfg.ChannelCount,
		Format:       cfg.Format,
		BufferSize:   time.Duration(cfg.BufferSize),
	}

	// Remember that you should **not** create more than one context
	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return fmt.Errorf("error: create oto context failed: %s", err)
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan
	p.ctx = ctx
	p.initialized = true
	return nil
}

func (p *player) Play(r io.ReadCloser) error {
	// Initialize the audio
	audio := p.ctx.NewPlayer(r)
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	audio.Play()

	// We can wait for the sound to finish playing using something like this
	for audio.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// If you don't want the audio/sound anymore simply close
	if err := audio.Close(); err != nil {
		return fmt.Errorf("error: close audio failed: %s", err)
	}

	return nil
}

func GetPlayer() Player {
	return p
}

func InitPlayer(cfg *config.Engine) {
	if err := p.init(cfg); err != nil {
		panic(fmt.Sprintf("error: init player failed: %s", err))
	}
}
