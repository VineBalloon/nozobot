package handlers

import "github.com/VineBalloon/nozobot/client"

// Stop
// The command to stop any audio streaming session gracefully
type Stop struct {
	Name string
}

func (p *Stop) Desc() string {
	return "Stops all audio in a pinchi :sparkling_heart:"
}

func (p *Stop) Roles() []string {
	return nil
}

func (p *Stop) Channels() []string {
	return nil
}

// Handle
// Tries to stop the current streaming session if there is one.
func (p *Stop) Handle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	// Stop the streaming session
	err := cs.StopStream()
	if err != nil {
		return err
	}

	// Signal that we've stopped
	_, err = s.ChannelMessageSend(m.ChannelID, "Stopped!")
	if err != nil {
		return err
	}

	return nil
}

func NewStop() *Stop {
	return &Stop{
		"Stop",
	}
}
