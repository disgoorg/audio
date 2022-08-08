package disgoplayer

import (
	"github.com/disgoorg/disgoplayer/channelconverter"
	"github.com/disgoorg/disgoplayer/opus"
	"github.com/disgoorg/snowflake/v2"
)

func NewPCMFrameChannelConverterCombinedReceiver(receiver PCMCombinedFrameReceiver, rate int, inputChannels int, outputChannels int) PCMCombinedFrameReceiver {
	return &pcmFrameChannelConverterCombinedReceiver{
		r:                receiver,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, opus.GetOutputBuffSize(rate, outputChannels)),
	}
}

type pcmFrameChannelConverterCombinedReceiver struct {
	r                PCMCombinedFrameReceiver
	channelConverter *channelconverter.ChannelConverter
	newPCM           []int16
}

func (p *pcmFrameChannelConverterCombinedReceiver) ReceiveCombinedPCMFrame(userIDs []snowflake.ID, packet *CombinedPCMPacket) error {
	if err := p.channelConverter.Convert(packet.PCM, p.newPCM); err != nil {
		return err
	}
	packet.PCM = p.newPCM
	return p.r.ReceiveCombinedPCMFrame(userIDs, packet)
}

func (*pcmFrameChannelConverterCombinedReceiver) Close() {}
