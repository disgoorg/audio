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

func (P *pcmVolumeFrameProvider) ProvidePCMFrame() ([]int16, error) {
	frame, err := P.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}
	applyVolume(frame, P.volumeProvider())
	return frame, nil
}

func (P *pcmVolumeFrameProvider) Close() {
	P.pcmFrameProvider.Close()
}

func applyVolume(in []int16, newVolume float32) {
	if newVolume == 1 {
		return
	}
	for i := range in {
		if newVolume == 0 {
			in[i] = 0
			continue
		}
		v := float32(in[i]) * newVolume
		if v > 32767 {
			v = 32767
		}
		if v < -32768 {
			v = -32768
		}
		in[i] = int16(v)
	}
}
