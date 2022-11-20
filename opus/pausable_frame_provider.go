package opus

import (
	"github.com/disgoorg/disgo/voice"
)

func NewPauseableFrameProvider(provider voice.OpusFrameProvider, pausedFunc func() bool) voice.OpusFrameProvider {
	return &pauseableFrameProvider{
		provider:   provider,
		pausedFunc: pausedFunc,
	}
}

type pauseableFrameProvider struct {
	provider   voice.OpusFrameProvider
	pausedFunc func() bool
}

func (p *pauseableFrameProvider) ProvideOpusFrame() ([]byte, error) {
	if p.pausedFunc() {
		return nil, nil
	}
	return p.provider.ProvideOpusFrame()
}

func (p *pauseableFrameProvider) Close() {
	p.provider.Close()
}
