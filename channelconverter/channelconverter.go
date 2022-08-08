package channelconverter

import "errors"

var ErrInvalidChannelCount = errors.New("invalid channel count")

func CreateChannelConverter(inputChannels int, outputChannels int) *ChannelConverter {
	return &ChannelConverter{
		inputChannels:  inputChannels,
		outputChannels: outputChannels,
	}
}

type ChannelConverter struct {
	inputChannels  int
	outputChannels int
}

func (c *ChannelConverter) Convert(input []int16, output []int16) error {
	for i := 0; i < len(input); i += c.inputChannels {
		if c.inputChannels == 1 && c.outputChannels == 2 {
			output[i*2] = input[i]
			output[i*2+1] = input[i]
		} else if c.inputChannels == 2 && c.outputChannels == 1 {
			newOutput := (int32(input[i]) + int32(input[i+1])) / 2
			if newOutput > 32767 {
				newOutput = 32767
			}
			if newOutput < -32768 {
				newOutput = -32768
			}
			output[i/2] = int16(newOutput)
		} else {
			return ErrInvalidChannelCount
		}
	}
	return nil
}
func (c *ChannelConverter) InputChannels() int {
	return c.inputChannels
}

func (c *ChannelConverter) OutputChannels() int {
	return c.inputChannels
}
