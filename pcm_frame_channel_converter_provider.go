package disgoplayer

import "github.com/disgoorg/disgoplayer/channelconverter"

func NewPCMFrameChannelConverterProvider(Provider PCMFrameProvider, inputChannels int, outputChannels int) PCMFrameProvider {
	return &pcmFrameChannelConverterProvider{
		pcmFrameProvider: Provider,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, outputChannels),
	}
}

type pcmFrameChannelConverterProvider struct {
	pcmFrameProvider PCMFrameProvider
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
