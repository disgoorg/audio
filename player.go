package audio

import (
	"io"
	"sync"

	"github.com/disgoorg/audio/pcm"
	"github.com/disgoorg/disgo/voice"
)

type Player interface {
	voice.OpusFrameProvider

	Volume() float32
	SetVolume(volume float32)
	Paused() bool
	SetPaused(paused bool)
}

func NewPlayer(providerFunc func() pcm.FrameProvider, listeners ...Listener) (Player, error) {
	player := &defaultPlayer{
		listeners: listeners,
		volume:    1,
		paused:    false,
	}

	pauseableProvider := pcm.NewPauseablePCMFrameProvider(pcm.NewVariablePCMFrameProvider(providerFunc), func() bool {
		return player.paused
	})

	volumeProvider := pcm.NewPCMVolumeFrameProvider(pauseableProvider, func() float32 {
		return player.volume
	})

	var err error
	if player.opusFrameProvider, err = pcm.NewPCMOpusProvider(nil, volumeProvider); err != nil {
		return nil, err
	}

	return player, nil
}

type defaultPlayer struct {
	opusFrameProvider voice.OpusFrameProvider
	volume            float32
	paused            bool
	playing           bool
	mu                sync.Mutex

	listeners []Listener
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
	if p.paused == paused {
		p.mu.Unlock()
		return
	}
	p.paused = paused
	p.mu.Unlock()
	if paused {
		p.emit(func(l Listener) {
			l.OnPause(p)
		})
	} else {
		p.emit(func(l Listener) {
			l.OnResume(p)
		})
	}
}

func (p *defaultPlayer) ProvideOpusFrame() ([]byte, error) {
	frame, err := p.opusFrameProvider.ProvideOpusFrame()
	if err == io.EOF {
		p.playing = false
		p.emit(func(l Listener) {
			l.OnEnd(p)
		})
	} else if err != nil {
		p.emit(func(l Listener) {
			l.OnError(p, err)
		})
	}
	if frame != nil && !p.playing {
		p.playing = true
		p.emit(func(l Listener) {
			l.OnStart(p)
		})
	}
	return frame, err
}

func (p *defaultPlayer) Close() {
	p.opusFrameProvider.Close()
	p.emit(func(l Listener) {
		l.OnClose(p)
	})
}

func (p *defaultPlayer) emit(l func(l Listener)) {
	for _, listener := range p.listeners {
		l(listener)
	}
}

type Listener interface {
	OnPause(player Player)
	OnResume(player Player)
	OnStart(player Player)
	OnEnd(player Player)
	OnError(player Player, err error)
	OnClose(player Player)
}
