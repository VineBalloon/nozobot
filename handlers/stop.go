package handlers

import (
	"errors"

	"github.com/VineBalloon/nozobot/client"
)

// Stop
// The command to stop any audio streaming session gracefully
type Stop struct {
	name        string
	description string
}

func (s *Stop) Name() string {
	return s.name
}

func (s *Stop) Desc() string {
	return s.description
}

func (s *Stop) Roles() []string {
	return nil
}

func (s *Stop) Channels() []string {
	return nil
}

// Handle
// Tries to stop the current streaming session if there is one.
func (s *Stop) MsgHandle(cs *client.ClientState) error {
	ss := cs.Session
	m := cs.Message

	// Check for voice channel
	if cs.Voice == nil {
		return errors.New("stop: no voice room to stop")
	}

	// Stop the streaming session
	err := cs.Voice.Stop()
	if err != nil {
		return err
	}

	// Signal that we've stopped
	_, err = ss.ChannelMessageSend(m.ChannelID, "Stopped!")
	if err != nil {
		return err
	}

	return nil
}

func NewStop() *Stop {
	return &Stop{
		"Stop",
		"Stops all audio in a pinchi :sparkling_heart:",
	}
}
