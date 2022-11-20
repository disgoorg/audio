package pcm

func NewPauseablePCMFrameProvider(provider FrameProvider, pauseProvider func() bool) FrameProvider {
	return &pauseablePCMFrameProvider{
		provider:      provider,
		pauseProvider: pauseProvider,
	}
}

type pauseablePCMFrameProvider struct {
	provider      FrameProvider
	pauseProvider func() bool
}

func (p *pauseablePCMFrameProvider) ProvidePCMFrame() ([]int16, error) {
	if p.pauseProvider() {
		return nil, nil
	}
	return p.provider.ProvidePCMFrame()
}

func (p *pauseablePCMFrameProvider) Close() {
	p.provider.Close()
}
