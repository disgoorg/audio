package samplerate

import (
	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/audio/pcm"
	"github.com/disgoorg/snowflake/v2"
)

// NewPCMFrameResamplerReceiver creates a FrameReceiver that resamples the PCM frames to the specified sample rate.
// If the resampler is nil, it will be created with samplerate.ConverterTypeSincBestQuality.
// The input sample rate is the sample rate of the PCM frames provided by the receiver.
// The output sample rate is the sample rate of the PCM frames returned by the receiver.
// The channels are the number of channels of the PCM frames provided by the receiver.
func NewPCMFrameResamplerReceiver(resampler *Resampler, inputSampleRate int, outputSampleRate int, channels int, pcmFrameReceiver pcm.FrameReceiver) pcm.FrameReceiver {
	if resampler == nil {
		resampler = CreateResampler(ConverterTypeSincBestQuality, channels)
	}
	return &sampleRateReceiver{
		resampler:        resampler,
		pcmFrameReceiver: pcmFrameReceiver,
		inputSampleRate:  inputSampleRate,
		outputSampleRate: outputSampleRate,
		newPCM:           make([]int16, opus.GetOutputBuffSize(outputSampleRate, channels)),
	}
}

type sampleRateReceiver struct {
	resampler        *Resampler
	pcmFrameReceiver pcm.FrameReceiver
	inputSampleRate  int
	outputSampleRate int
	newPCM           []int16
}

func (p *sampleRateReceiver) ReceivePCMFrame(userID snowflake.ID, packet *pcm.Packet) error {

	var (
		inputFrames  int64
		outputFrames int64
	)
	if err := p.resampler.Process(packet.PCM, p.newPCM, p.inputSampleRate, p.outputSampleRate, 0, &inputFrames, &outputFrames); err != nil {
		return err
	}

	packet.PCM = p.newPCM
	return p.pcmFrameReceiver.ReceivePCMFrame(userID, packet)
}

func (p *sampleRateReceiver) CleanupUser(userID snowflake.ID) {
	p.pcmFrameReceiver.CleanupUser(userID)
}

func (p *sampleRateReceiver) Close() {
	p.resampler.Destroy()
	p.pcmFrameReceiver.Close()
}
