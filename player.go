package disgoplayer

import (
	"sync"

	"github.com/disgoorg/disgo/voice"
)

type Player interface {
	voice.OpusFrameProvider

	Volume() float32
	SetVolume(volume float32)
	Paused() bool
	SetPaused(paused bool)
}

func NewPlayer(provider PCMFrameProvider) (Player, error) {
	player := &defaultPlayer{
		provider: provider,
		volume:   1,
		paused:   false,
	}

	player.volumePCMProvider = NewPCMVolumeFrameProvider(provider, func() float32 {
		return player.volume
	})

	var err error
	if player.opusFrameProvider, err = NewPCMOpusProvider(nil, player); err != nil {
		return nil, err
	}

	return player, nil
}

type defaultPlayer struct {
	provider          PCMFrameProvider
	volumePCMProvider PCMFrameProvider
	opusFrameProvider voice.OpusFrameProvider
	volume            float32
	paused            bool
	mu                sync.Mutex
}

func (p *defaultPlayer) Volume() float32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.volume
}

func (p *defaultPlayer) SetVolume(volume float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.volume = volume
}

func (p *defaultPlayer) Paused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

func (p *defaultPlayer) SetPaused(paused bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = paused
}

func (p *defaultPlayer) ProvidePCMFrame() ([]int16, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.paused {
		return nil, nil
	}

	if p.volume != 1 {
		return p.volumePCMProvider.ProvidePCMFrame()
	}
	return p.provider.ProvidePCMFrame()
}

func (p *defaultPlayer) ProvideOpusFrame() ([]byte, error) {
	return p.opusFrameProvider.ProvideOpusFrame()
}

func (p *defaultPlayer) Close() {
	p.provider.Close()
}
