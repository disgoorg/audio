package disgoplayer

func NewPCMVolumeFrameProvider(pcmFrameProvider PCMFrameProvider, volumeProvider func() float32) PCMFrameProvider {
	return &pcmVolumeFrameProvider{
		pcmFrameProvider: pcmFrameProvider,
		volumeProvider:   volumeProvider,
	}
}

type pcmVolumeFrameProvider struct {
	pcmFrameProvider PCMFrameProvider
	volumeProvider   func() float32
}

func (p *pcmVolumeFrameProvider) ProvidePCMFrame() ([]int16, error) {
	frame, err := p.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}
	applyVolume(frame, p.volumeProvider())
	return frame, nil
}

func (p *pcmVolumeFrameProvider) Close() {
	p.pcmFrameProvider.Close()
}

func applyVolume(pcm []int16, newVolume float32) {
	if newVolume == 1 {
		return
	}
	for i := range pcm {
		if newVolume == 0 {
			pcm[i] = 0
			continue
		}
		v := float32(pcm[i]) * newVolume
		if v > 32767 {
			v = 32767
		}
		if v < -32768 {
			v = -32768
		}
		pcm[i] = int16(v)
	}
}
