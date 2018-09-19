package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	//"github.com/bwmarrin/discordgo"
)

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

func (p *Ping) Handle(c *client.ClientState) error {
	s := c.Session
	m := c.Message

	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}

func NewPing(n string) *Ping {
	return &Ping{
		"Ping",
	}
}
