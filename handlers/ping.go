package handlers

import (
	"github.com/VineBalloon/nozobot/client"
)

// Ping
// The command to ping the user when they ask for it
type Ping struct {
	Name string
}

func (p *Ping) Desc() string {
	return "Ping pong with Nozomi :ping_pong:"
}

func (p *Ping) Roles() []string {
	return nil
}

func (p *Ping) Channels() []string {
	return nil
}

// Handle
// Responds with "Pong!" when user sends a ping command
func (p *Ping) Handle(cs *client.ClientState) error {
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
	}
}
