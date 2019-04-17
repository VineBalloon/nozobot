package handlers

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"

	"github.com/VineBalloon/nozobot/client"
)

var (
	MAJOR_DIR string = "./tarots/major-arcana/"
	MINOR_DIR string = "./tarots/minor-arcana/"
)

// Tarot
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

	// Read from major arcana, cwd is nozobot's root
	majors, err := ioutil.ReadDir(MAJOR_DIR)
	if err != nil {
		return err
	}
	major := majors[rand.Int()%len(majors)]
	path := MAJOR_DIR + major.Name()
	ext := filepath.Ext(path)[1:]
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	img := &discordgo.File{
		major.Name(),
		"image/" + ext,
		fd,
	}
	msg := "Nozomi Spiritual Power!"
	out := &discordgo.MessageSend{
		msg,
		nil,
		false,
		[]*discordgo.File{img},
		nil,
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
