package pcm

import (
	"fmt"
	"sync"

	"github.com/disgoorg/audio/opus"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

// NewPCMOpusReceiver creates a new voice.OpusFrameReceiver which receives Opus frames and decodes them into PCM frames. A new decoder is created for each user.
// You can pass your own *opus.Decoder by passing a decoderCreateFunc or nil to use the default Opus decoder(48000hz sample rate, 2 channels).
// You can filter users by passing a voice.ShouldReceiveUserFunc or nil to receive all users.
func NewPCMOpusReceiver(decoderCreateFunc func() (*opus.Decoder, error), pcmFrameReceiver FrameReceiver, userFilter voice.UserFilterFunc) voice.OpusFrameReceiver {
	if decoderCreateFunc == nil {
		decoderCreateFunc = func() (*opus.Decoder, error) {
			decoder, err := opus.NewDecoder(48000, 2)
			if err != nil {
				return nil, fmt.Errorf("failed to create opus decoder: %w", err)
			}
			return decoder, nil
		}
	}
	return &pcmOpusReceiver{
		userFilter:        userFilter,
		decoderCreateFunc: decoderCreateFunc,
		decoderStates:     map[snowflake.ID]*decoderState{},
		pcmFrameReceiver:  pcmFrameReceiver,
	}
}

type decoderState struct {
	decoder *opus.Decoder
	pcmBuff []int16
}

type pcmOpusReceiver struct {
	userFilter        voice.UserFilterFunc
	decoderCreateFunc func() (*opus.Decoder, error)
	decoderStates     map[snowflake.ID]*decoderState
	decodersMu        sync.Mutex
	pcmFrameReceiver  FrameReceiver
}

func (r *pcmOpusReceiver) ReceiveOpusFrame(userID snowflake.ID, packet *voice.Packet) error {
	if r.userFilter != nil && !r.userFilter(userID) {
		return nil
	}
	r.decodersMu.Lock()
	state, ok := r.decoderStates[userID]
	if !ok {
		decoder, err := r.decoderCreateFunc()
		if err != nil {
			r.decodersMu.Unlock()
			return fmt.Errorf("failed to create opus decoder: %w", err)
		}

		sampleRate, err := decoder.SampleRate()
		if err != nil {
			r.decodersMu.Unlock()
			return fmt.Errorf("failed to get sample rate: %w", err)
		}

		state = &decoderState{
			decoder: decoder,
			pcmBuff: make([]int16, opus.GetOutputBuffSize(sampleRate, decoder.Channels())),
		}
		r.decoderStates[userID] = state
	}
	r.decodersMu.Unlock()

	_, err := state.decoder.Decode(packet.Opus, state.pcmBuff, false)
	if err != nil {
		return err
	}

	return r.pcmFrameReceiver.ReceivePCMFrame(userID, &Packet{
		SSRC:      packet.SSRC,
		Sequence:  packet.Sequence,
		Timestamp: packet.Timestamp,
		PCM:       state.pcmBuff,
	})
}

func (r *pcmOpusReceiver) CleanupUser(userID snowflake.ID) {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	state, ok := r.decoderStates[userID]
	if ok {
		state.decoder.Destroy()
		delete(r.decoderStates, userID)
	}
	r.pcmFrameReceiver.CleanupUser(userID)
}

func (r *pcmOpusReceiver) Close() {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	for _, state := range r.decoderStates {
		state.decoder.Destroy()
	}
	r.pcmFrameReceiver.Close()
}
