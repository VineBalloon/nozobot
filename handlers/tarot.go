package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/bwmarrin/discordgo"
)

type Tarot struct {
	Name string
}

func (t *Tarot) Desc() string {
	return "Nozomi decides your fate!"
}

func (t *Tarot) Roles() []string {
	return nil
}

func (t *Tarot) Channels() []string {
	return nil
}

func (t *Tarot) HandleTarot(c *client.ClientState) error {
	s := c.Session
	m := c.Message

	return nil
}

func NewTarot(n string) *Tarot {
	return &Tarot{
		"Tarot",
	}
}
