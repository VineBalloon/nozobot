package handlers

import (
	"errors"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/helpers"
)

// Leave
// The command to leave voice gracefully
type Leave struct {
	name        string
	description string
}

func (p *Leave) Name() string {
	return p.name
}

func (p *Leave) Desc() string {
	return p.description
}

func (p *Leave) Roles() []string {
	return nil
}

func (p *Leave) Channels() []string {
	return nil
}

// Handle
// Gracefully closes all encoding, stream, and voice connections
func (p *Leave) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	// Check voice channel
	if cs.Voice == nil {
		return errors.New("leave: no voice room to leave")
	}

	// Get voice channel info
	channel, err := cs.Session.Channel(cs.Voice.Id)
	if err != nil {
		return err
	}

	// Leave the voice connection
	//cs.Voice.Stop()
	//cs.Voice.Connection.Close()
	err = cs.Voice.Leave()
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
		"Leaves all audio in a pinchi :sparkling_heart:",
	}
}
