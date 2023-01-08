package pcm

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
)

func NewBufferPCMProvider(provider FrameProvider) FrameProvider {
	ctx, cancel := context.WithCancel(context.Background())
	buffer := &bufferFrameProvider{
		provider: provider,
		cancel:   cancel,
	}

	go buffer.process(ctx)
	return buffer
}

type bufferFrameProvider struct {
	provider   FrameProvider
	buff       [][]int16
	processing atomic.Bool
	cancel     context.CancelFunc
}

func (p *bufferFrameProvider) process(ctx context.Context) {
	defer p.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(p.buff) >= 10 {
				continue
			}
			frame, err := p.provider.ProvidePCMFrame()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				continue
			}
			p.buff = append(p.buff, frame)
		}
	}
}

func (p *bufferFrameProvider) ProvidePCMFrame() (frame []int16, err error) {
	if len(p.buff) == 0 {
		return
	}
	frame, p.buff = p.buff[0], p.buff[1:]
	return
}

func (p *bufferFrameProvider) Close() {
	p.cancel()
	p.provider.Close()
}
