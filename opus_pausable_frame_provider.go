package disgoplayer

import (
	"github.com/disgoorg/disgo/voice"
)

func NewPauseableOpusFrameProvider(opusProvider voice.OpusFrameProvider, pauseProvider func() bool) voice.OpusFrameProvider {
	return &pauseableOpusFrameProvider{
		opusProvider:  opusProvider,
		pauseProvider: pauseProvider,
	}
}

type pauseableOpusFrameProvider struct {
	opusProvider  voice.OpusFrameProvider
	pauseProvider func() bool
}

func (p *pauseableOpusFrameProvider) ProvideOpusFrame() ([]byte, error) {
	if p.pauseProvider() {
		return nil, nil
	}
	return p.opusProvider.ProvideOpusFrame()
}

func (p *pauseableOpusFrameProvider) Close() {
	p.opusProvider.Close()
}
