package handlers

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/VineBalloon/nozobot/client"
)

// Imit
// The command to convert a string into a randomly upper-lower version
type Imit struct {
	name        string
	description string
}

func (i *Imit) Name() string {
	return i.name
}

func (i *Imit) Desc() string {
	return i.description
}

func (i *Imit) Roles() []string {
	return []string{"gay", "boi"}
}

func (i *Imit) Channels() []string {
	return nil
}

// MsgHandle
func (g *Imit) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	in := strings.Join(cs.Arguments, " ")
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

func NewImit() *Imit {
	return &Imit{
		"Imit",
		"Convert a string into a randomly upper/lower-cased version",
	}
}
