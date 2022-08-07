package disgoplayer

import (
	"github.com/disgoorg/disgoplayer/channelconverter"
	"github.com/disgoorg/snowflake/v2"
)

func NewPCMFrameChannelConverterReceiver(receiver PCMFrameReceiver, inputChannels int, outputChannels int) PCMFrameReceiver {
	return &pcmFrameChannelConverterReceiver{
		r:                receiver,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, outputChannels),
	}
}

type pcmFrameChannelConverterReceiver struct {
	r                PCMFrameReceiver
	channelConverter *channelconverter.ChannelConverter
	newPCM           []int16
}

func (p *pcmFrameChannelConverterReceiver) ReceivePCMFrame(userID snowflake.ID, packet *PCMPacket) error {
	if err := p.channelConverter.Convert(packet.PCM, p.newPCM); err != nil {
		return err
	}
	packet.PCM = p.newPCM
	return p.r.ReceivePCMFrame(userID, packet)
}

func (p *pcmFrameChannelConverterReceiver) CleanupUser(_ snowflake.ID) {}

func (*pcmFrameChannelConverterReceiver) Close() {}
