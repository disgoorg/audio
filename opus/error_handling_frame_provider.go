package opus

import (
	"errors"

	"github.com/disgoorg/disgo/voice"
)

// ErrProviderClosed indicates that the voice.OpusFrameProvider was closed.
var ErrProviderClosed = errors.New("opus provider closed")

// ErrorHandleFunc is a function which is called when an error is returned by the given voice.OpusFrameProvider.
type ErrorHandleFunc func(err error)

// NewErrorHandlingFrameProvider creates a new voice.OpusFrameProvider which intercepts errors returned by the given voice.OpusFrameProvider and calls the given ErrorHandleFunc.
// This can be used to intercept io.EOF or other errors and stop the playback or start a new track.
func NewErrorHandlingFrameProvider(opusFrameProvider voice.OpusFrameProvider, errorHandleFunc ErrorHandleFunc) voice.OpusFrameProvider {
	return &errorHandlingFrameProvider{
		opusFrameProvider: opusFrameProvider,
		errorHandleFunc:   errorHandleFunc,
	}
}

type errorHandlingFrameProvider struct {
	opusFrameProvider voice.OpusFrameProvider
	errorHandleFunc   ErrorHandleFunc
}

func (p *errorHandlingFrameProvider) ProvideOpusFrame() ([]byte, error) {
	frame, err := p.opusFrameProvider.ProvideOpusFrame()
	if err != nil {
		p.errorHandleFunc(err)
		return nil, err
	}
	return frame, nil
}

func (p *errorHandlingFrameProvider) Close() {
	p.opusFrameProvider.Close()
	p.errorHandleFunc(ErrProviderClosed)
}
