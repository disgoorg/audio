package pcm

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/disgoorg/audio/opus"
)

// FrameProvider is an interface for providing PCM frames.
type FrameProvider interface {
	// ProvidePCMFrame is called to get a PCM frame.
	ProvidePCMFrame() ([]int16, error)

	// Close is called when the provider is no longer needed. It should close any open resources.
	Close()
}

// NewReader creates a new FrameProvider which reads PCM frames from the given io.Reader.
func NewReader(r io.Reader) FrameProvider {
	return NewCustomReader(r, 48000, 2)
}

// NewCustomReader creates a new FrameProvider which reads PCM frames from the given io.Reader.
// You can specify the sample rate and number of channels.
func NewCustomReader(r io.Reader, rate int, channels int) FrameProvider {
	return &reader{
		r:           r,
		bytePCMBuff: make([]byte, opus.GetOutputBuffSize(rate, channels)*2),
		pcmBuff:     make([]int16, opus.GetOutputBuffSize(rate, channels)),
	}
}

type reader struct {
	r           io.Reader
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *reader) ProvidePCMFrame() ([]int16, error) {
	_, err := p.r.Read(p.bytePCMBuff)
	if err != nil {
		return nil, err
	}

	if err = binary.Read(bytes.NewReader(p.bytePCMBuff), binary.LittleEndian, p.pcmBuff); err != nil {
		return nil, err
	}
	return p.pcmBuff, nil
}

func (*reader) Close() {}
