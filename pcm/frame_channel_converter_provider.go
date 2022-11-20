package pcm

import (
	"github.com/disgoorg/audio/channelconverter"
	"github.com/disgoorg/audio/opus"
)

func NewPCMFrameChannelConverterProvider(Provider FrameProvider, rate int, inputChannels int, outputChannels int) FrameProvider {
	return &pcmFrameChannelConverterProvider{
		pcmFrameProvider: Provider,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, opus.GetOutputBuffSize(rate, outputChannels)),
	}
}

type pcmFrameChannelConverterProvider struct {
	pcmFrameProvider FrameProvider
	channelConverter *channelconverter.ChannelConverter
	newPCM           []int16
}

func (p *pcmFrameChannelConverterProvider) ProvidePCMFrame() ([]int16, error) {
	frame, err := p.pcmFrameProvider.ProvidePCMFrame()
	if err != nil {
		return nil, err
	}

	if err = p.channelConverter.Convert(frame, p.newPCM); err != nil {
		return nil, err
	}
	return p.newPCM, nil
}

func (*pcmFrameChannelConverterProvider) Close() {}
