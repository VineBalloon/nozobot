package handlers

import (
	"github.com/VineBalloon/nozobot/client"
)

// TODO Tarot
// The command to generate random fortunes
type Tarot struct {
	name        string
	description string
}

func (t *Tarot) Name() string {
	return t.name
}

func (t *Tarot) Desc() string {
	return t.description
}

func (t *Tarot) Roles() []string {
	return nil
}

func (t *Tarot) Channels() []string {
	return nil
}

func (t *Tarot) MsgHandle(cs *client.ClientState) error {
	//s := cs.Session
	//m := cs.Message

	return nil
}

func NewTarot() *Tarot {
	return &Tarot{
		"Tarot",
		"Nozomi decides your fate!",
	}
}
