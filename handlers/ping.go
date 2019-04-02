package handlers

import (
	"github.com/VineBalloon/nozobot/client"
)

// Ping The command to ping the user when they ask for it
type Ping struct {
	name        string
	description string
}

// Name Returns name of the command
func (p *Ping) Name() string {
	return p.name
}

// Desc Returns description of the command
func (p *Ping) Desc() string {
	return p.description
}

// Roles Returns roles required by the command
func (p *Ping) Roles() []string {
	return nil
}

// Channels Returns channels required by the command
func (p *Ping) Channels() []string {
	return nil
}

// MsgHandle Responds with "Pong!" when user sends a ping command
func (p *Ping) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}

// NewPing Constructor for Ping
func NewPing() *Ping {
	return &Ping{
		"Ping",
		"Ping pong with Nozomi :ping_pong:",
	}
}
