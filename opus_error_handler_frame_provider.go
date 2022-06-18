package disgoplayer

import (
	"errors"

	"github.com/disgoorg/disgo/voice"
)

// OpusProviderClosed indicates that the voice.OpusFrameProvider was closed.
var OpusProviderClosed = errors.New("opus provider closed")

// ErrorHandleFunc is a function which is called when an error is returned by the given voice.OpusFrameProvider.
type ErrorHandleFunc func(err error)

// NewErrorHandlerOpusFrameProvider creates a new voice.OpusFrameProvider which intercepts errors returned by the given voice.OpusFrameProvider and calls the given ErrorHandleFunc.
// This can be used to intercept io.EOF or other errors and stop the playback or start a new track.
func NewErrorHandlerOpusFrameProvider(opusFrameProvider voice.OpusFrameProvider, errorHandleFunc ErrorHandleFunc) voice.OpusFrameProvider {
	return &playerOpusFrameProvider{
		opusFrameProvider: opusFrameProvider,
		errorHandleFunc:   errorHandleFunc,
	}
}

type playerOpusFrameProvider struct {
	opusFrameProvider voice.OpusFrameProvider
	errorHandleFunc   ErrorHandleFunc
}

func (p *playerOpusFrameProvider) ProvideOpusFrame() ([]byte, error) {
	frame, err := p.opusFrameProvider.ProvideOpusFrame()
	if err != nil {
		p.errorHandleFunc(err)
		return nil, err
	}
	return frame, nil
}

func (p *playerOpusFrameProvider) Close() {
	p.opusFrameProvider.Close()
	p.errorHandleFunc(OpusProviderClosed)
}
