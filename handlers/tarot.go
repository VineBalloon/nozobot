package handlers

import (
	//"github.com/VineBalloon/nozobot/helpers"
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

func (t *Tarot) HandleTarot(s *discordgo.Session, m *discordgo.MessageCreate) error {
	return nil
}

func NewTarot(n string) *Tarot {
	return &Tarot{
		"Tarot",
	}
}
