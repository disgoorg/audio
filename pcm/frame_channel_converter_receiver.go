package pcm

import (
	"github.com/disgoorg/audio/channelconverter"
	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/snowflake/v2"
)

func NewPCMFrameChannelConverterReceiver(receiver FrameReceiver, rate int, inputChannels int, outputChannels int) FrameReceiver {
	return &pcmFrameChannelConverterReceiver{
		r:                receiver,
		channelConverter: channelconverter.CreateChannelConverter(inputChannels, outputChannels),
		newPCM:           make([]int16, opus.GetOutputBuffSize(rate, outputChannels)),
	}
}

type pcmFrameChannelConverterReceiver struct {
	r                FrameReceiver
	channelConverter *channelconverter.ChannelConverter
	newPCM           []int16
}

func (p *pcmFrameChannelConverterReceiver) ReceivePCMFrame(userID snowflake.ID, packet *Packet) error {
	if err := p.channelConverter.Convert(packet.PCM, p.newPCM); err != nil {
		return err
	}
	packet.PCM = p.newPCM
	return p.r.ReceivePCMFrame(userID, packet)
}

func (p *pcmFrameChannelConverterReceiver) CleanupUser(_ snowflake.ID) {}

func (*pcmFrameChannelConverterReceiver) Close() {}
