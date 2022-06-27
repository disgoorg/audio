package disgoplayer

import (
	"github.com/disgoorg/disgo/voice"
)

func NewPauseableOpusFrameProvider(provider voice.OpusFrameProvider, pausedFunc func() bool) voice.OpusFrameProvider {
	return &pauseableOpusFrameProvider{
		provider:   provider,
		pausedFunc: pausedFunc,
	}
}

type pauseableOpusFrameProvider struct {
	provider   voice.OpusFrameProvider
	pausedFunc func() bool
}

func (p *pauseableOpusFrameProvider) ProvideOpusFrame() ([]byte, error) {
	if p.pausedFunc() {
		return nil, nil
	}
	return p.provider.ProvideOpusFrame()
}

func (p *pauseableOpusFrameProvider) Close() {
	p.provider.Close()
}
