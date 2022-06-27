package disgoplayer

func NewPauseablePCMFrameProvider(provider PCMFrameProvider, pauseProvider func() bool) PCMFrameProvider {
	return &pauseablePCMFrameProvider{
		provider:      provider,
		pauseProvider: pauseProvider,
	}
}

type pauseablePCMFrameProvider struct {
	provider      PCMFrameProvider
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
