package samplerate

import (
	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/audio/pcm"
)

// NewPCMFrameResamplerProvider creates a FrameProvider that resamples the PCM frames to the specified sample rate.
// If the resampler is nil, it will be created with samplerate.ConverterTypeSincBestQuality.
// The input sample rate is the sample rate of the PCM frames provided by the provider.
// The output sample rate is the sample rate of the PCM frames returned by the provider.
// The channels are the number of channels of the PCM frames provided by the provider.
func NewPCMFrameResamplerProvider(resampler *Resampler, inputSampleRate int, outputSampleRate int, channels int, pcmFrameProvider pcm.FrameProvider) pcm.FrameProvider {
	if resampler == nil {
		resampler = CreateResampler(ConverterTypeSincBestQuality, channels)
	}
	return &sampleRateProvider{
		resampler:        resampler,
		pcmFrameProvider: pcmFrameProvider,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
		newPCM:           make([]int16, opus.GetOutputBuffSize(outputSampleRate, channels)),
	}
}

type sampleRateProvider struct {
	resampler        *Resampler
	pcmFrameProvider pcm.FrameProvider
	inputSampleRate  int
	outputSampleRate int
	newPCM           []int16
}

func (p *sampleRateProvider) ProvidePCMFrame() ([]int16, error) {
	pcm, err := p.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}

	var (
		inputFrames  int64
		outputFrames int64
	)
	if err = p.resampler.Process(pcm, p.newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return nil, err
	}

	return p.newPCM, nil
}

func (p *sampleRateProvider) Close() {
	p.resampler.Destroy()
	p.pcmFrameProvider.Close()
}
