package pcm

import (
	"encoding/binary"
	"io"

	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

type (
	// FrameReceiver is an interface for receiving PCM frames.
	FrameReceiver interface {
		// ReceivePCMFrame is called when a PCM frame is received.
		ReceivePCMFrame(userID snowflake.ID, packet *Packet) error

		// CleanupUser is called when a user is disconnected. This should close any resources associated with the user.
		CleanupUser(userID snowflake.ID)

		// Close is called when the receiver is no longer needed. It should close any open resources.
		Close()
	}

	// Packet is a 20ms PCM frame with a ssrc, sequence and timestamp.
	Packet struct {
		SSRC      uint32
		Sequence  uint16
		Timestamp uint32
		PCM       []int16
	}
)

// NewWriter creates a new FrameReceiver which writes PCM frames to the given io.Writer.
// You can filter which users should be written by passing a voice.ShouldReceiveUserFunc.
func NewWriter(w io.Writer, userFilter voice.UserFilterFunc) FrameReceiver {
	return &writer{
		w:          w,
		userFilter: userFilter,
	}
}

type writer struct {
	w          io.Writer
	userFilter voice.UserFilterFunc
}

func (p *writer) ReceivePCMFrame(userID snowflake.ID, packet *Packet) error {
	if p.userFilter == nil && !p.userFilter(userID) {
		return nil
	}
	return binary.Write(p.w, binary.LittleEndian, packet.PCM)
}

func (p *writer) CleanupUser(_ snowflake.ID) {}

func (*writer) Close() {}
