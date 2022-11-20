package pcm

func NewVariablePCMFrameProvider(providerFunc func() FrameProvider) FrameProvider {
	return &variablePCMFrameProvider{
		providerFunc: providerFunc,
	}
}

type variablePCMFrameProvider struct {
	providerFunc func() FrameProvider
}

func (v *variablePCMFrameProvider) ProvidePCMFrame() ([]int16, error) {
	if provider := v.providerFunc(); provider != nil {
		return provider.ProvidePCMFrame()
	}
	return nil, nil
}

func (v *variablePCMFrameProvider) Close() {
	if provider := v.providerFunc(); provider != nil {
		provider.Close()
	}
}
