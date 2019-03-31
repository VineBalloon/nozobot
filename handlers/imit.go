package handlers

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/VineBalloon/nozobot/client"
)

// Imit The command to convert a string into a randomly upper-lower version
type Imit struct {
	name        string
	description string
}

// Name Returns name of the command
func (i *Imit) Name() string {
	return i.name
}

// Desc Returns description of the command
func (i *Imit) Desc() string {
	return i.description
}

// Roles Returns roles required by the command
func (i *Imit) Roles() []string {
	return []string{"gay", "boi"}
}

// Channels Returns channels required by the command
func (i *Imit) Channels() []string {
	return nil
}

// MsgHandle Convert a string into a randomly upper/lower-cased version
func (i *Imit) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	in := strings.Join(cs.Args, " ")
	out := []rune{}
	for _, c := range in {
		if rand.Intn(10)%2 == 0 {
			out = append(out, unicode.ToUpper(c))
		} else {
			out = append(out, unicode.ToLower(c))
		}
	}
	s.ChannelMessageSend(m.ChannelID, string(out))
	return nil
}

// NewImit Constructor for Imit
func NewImit() *Imit {
	return &Imit{
		"Imit",
		"Imitates you",
	}
}
