package disgoplayer

func NewVariablePCMFrameProvider(providerFunc func() PCMFrameProvider) PCMFrameProvider {
	return &variablePCMFrameProvider{
		providerFunc: providerFunc,
	}
}

type variablePCMFrameProvider struct {
	providerFunc func() PCMFrameProvider
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
