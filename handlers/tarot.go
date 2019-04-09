package handlers

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"

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
	s := cs.Session
	m := cs.Message

	// TODO: Get random file from tarots
	name := "tarots/example.png"
	ext := filepath.Ext(name)[1:]
	fd, err := os.Open(name)
	if err != nil {
		return err
	}
	img := &discordgo.File{
		name,
		"image/" + ext,
		fd,
	}
	msg := "Nozomi Spiritual Power!"
	out := &discordgo.MessageSend{
		msg,
		nil,
		false,
		[]*discordgo.File{img},
	}
	_, err = s.ChannelMessageSendComplex(m.ChannelID, out)
	if err != nil {
		return err
	}
	return nil
}

func NewTarot() *Tarot {
	return &Tarot{
		"Tarot",
		"Nozomi decides your fate!",
	}
}
