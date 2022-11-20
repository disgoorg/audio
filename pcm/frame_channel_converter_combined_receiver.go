package pcm

import (
	"github.com/disgoorg/audio/channelconverter"
	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/snowflake/v2"
)

func NewFrameChannelConverterCombinedReceiver(receiver CombinedFrameReceiver, rate int, inputChannels int, outputChannels int) CombinedFrameReceiver {
	return &frameChannelConverterCombinedReceiver{
		r:                receiver,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, opus.GetOutputBuffSize(rate, outputChannels)),
	}
}

type frameChannelConverterCombinedReceiver struct {
	r                CombinedFrameReceiver
	channelConverter *channelconverter.ChannelConverter
	newPCM           []int16
}

func (p *frameChannelConverterCombinedReceiver) ReceiveCombinedPCMFrame(userIDs []snowflake.ID, packet *CombinedPacket) error {
	if err := p.channelConverter.Convert(packet.PCM, p.newPCM); err != nil {
		return err
	}
	packet.PCM = p.newPCM
	return p.r.ReceiveCombinedPCMFrame(userIDs, packet)
}

func (*frameChannelConverterCombinedReceiver) Close() {}
