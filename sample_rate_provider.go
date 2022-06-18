package disgoplayer

import (
	"github.com/disgoorg/disgoplayer/opus"
	"github.com/disgoorg/disgoplayer/samplerate"
)

// NewPCMFrameResamplerProvider creates a PCMFrameProvider that resamples the PCM frames to the specified sample rate.
// If the resampler is nil, it will be created with samplerate.ConverterTypeSincBestQuality.
// The input sample rate is the sample rate of the PCM frames provided by the provider.
// The output sample rate is the sample rate of the PCM frames returned by the provider.
// The channels are the number of channels of the PCM frames provided by the provider.
func NewPCMFrameResamplerProvider(resampler *samplerate.Resampler, inputSampleRate int, outputSampleRate int, channels int, pcmFrameProvider PCMFrameProvider) PCMFrameProvider {
	if resampler == nil {
		resampler = samplerate.CreateResampler(samplerate.ConverterTypeSincBestQuality, channels)
	}
	return &sampleRateProvider{
		resampler:        resampler,
		pcmFrameProvider: pcmFrameProvider,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
	}
}

type sampleRateProvider struct {
	resampler        *samplerate.Resampler
	pcmFrameProvider PCMFrameProvider
	inputSampleRate  int
	outputSampleRate int
}

func (p *sampleRateProvider) ProvidePCMFrame() ([]int16, error) {
	pcm, err := p.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}

	newPCM := make([]int16, opus.GetOutputBuffSize(p.outputSampleRate, p.resampler.Channels()))
	var (
		inputFrames  int64
		outputFrames int64
	)
	if err = p.resampler.Process(pcm, newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return nil, err
	}

	return newPCM, nil
}

func (p *sampleRateProvider) Close() {
	p.resampler.Destroy()
	p.pcmFrameProvider.Close()
}
