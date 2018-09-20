package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/helpers"
)

// Leave
// The command to leave voice gracefully
type Leave struct {
	Name string
}

func (p *Leave) Desc() string {
	return "Leaves all audio in a pinchi :sparkling_heart:"
}

func (p *Leave) Roles() []string {
	return nil
}

func (p *Leave) Channels() []string {
	return nil
}

// Handle
// Gracefully closes all encoding, stream, and voice connections
func (p *Leave) Handle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	channel, err := cs.Session.Channel(cs.Voice.Id)
	if err != nil {
		return err
	}

	// Close all connections
	err = cs.Close()
	if err != nil {
		return err
	}

	// Signal that we've left
	_, err = s.ChannelMessageSend(m.ChannelID, "Left "+helpers.Code(channel.Name))
	if err != nil {
		return err
	}

	return nil
}

func NewLeave() *Leave {
	return &Leave{
		"Leave",
	}
}
