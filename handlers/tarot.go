package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	//"github.com/bwmarrin/discordgo"
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

func (t *Tarot) Handle(cs *client.ClientState) error {
	//s := cs.Session
	//m := cs.Message

	return nil
}

func NewTarot() *Tarot {
	return &Tarot{
		"Tarot",
	}
}
