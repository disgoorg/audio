package mp3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/audio/pcm"
)

// NewPCMFrameProvider returns a FrameProvider that reads mp3 and converts it into pcm frames.
// Write the Mp3 data to the returned writer.
func NewPCMFrameProvider(decoder *Decoder) (pcm.FrameProvider, io.Writer, error) {
	return NewCustomPCMFrameProvider(decoder, 48000, 2)
}

// NewCustomPCMFrameProvider returns a FrameProvider that reads mp3 and converts it into pcm frames.
// You can specify the rate and channels of the output PCM frames.
func NewCustomPCMFrameProvider(decoder *Decoder, rate int, channels int) (pcm.FrameProvider, io.Writer, error) {
	if decoder == nil {
		var err error
		decoder, err = CreateDecoder()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create mp3 decoder: %w", err)
		}
	}

	if err := decoder.Param(ForceRate, rate, float64(rate)); err != nil {
		return nil, nil, fmt.Errorf("failed to set param: %w", err)
	}

	if err := decoder.OpenFeed(); err != nil {
		return nil, nil, fmt.Errorf("failed to open feed for mp3 decoder: %w", err)
	}

	writeFunc := writer(func(p []byte) (int, error) {
		return decoder.Write(p)
	})

	return &pcmFrameProvider{
		decoder:     decoder,
		bytePCMBuff: make([]byte, opus.GetOutputBuffSize(rate, channels)*2),
		pcmBuff:     make([]int16, opus.GetOutputBuffSize(rate, channels)),
	}, writeFunc, nil
}

type pcmFrameProvider struct {
	decoder     *Decoder
	bytePCMBuff []byte
	pcmBuff     []int16
}

func (p *pcmFrameProvider) ProvidePCMFrame() ([]int16, error) {
	_, err := p.decoder.Read(p.bytePCMBuff)
	if err != nil {
		return nil, err
	}

	if err = binary.Read(bytes.NewReader(p.bytePCMBuff), binary.LittleEndian, p.pcmBuff); err != nil {
		return nil, err
	}
	return p.pcmBuff, nil
}

func (p *pcmFrameProvider) Close() {
	_ = p.decoder.Close()
}

type writer func(p []byte) (int, error)

func (w writer) Write(p []byte) (int, error) {
	return w(p)
}
