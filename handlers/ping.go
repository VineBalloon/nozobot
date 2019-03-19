package handlers

import (
	"github.com/VineBalloon/nozobot/client"
)

// Ping
// The command to ping the user when they ask for it
type Ping struct {
	name        string
	description string
}

func (p *Ping) Name() string {
	return p.name
}

func (p *Ping) Desc() string {
	return p.description
}

func (p *Ping) Roles() []string {
	return nil
}

func (p *Ping) Channels() []string {
	return nil
}

// Handle
// Responds with "Pong!" when user sends a ping command
func (p *Ping) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}

func NewPing() *Ping {
	return &Ping{
		"Ping",
		"Ping pong with Nozomi :ping_pong:",
	}
}
