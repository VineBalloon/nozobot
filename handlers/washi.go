package handlers

import (
	"time"

	"github.com/VineBalloon/nozobot/client"
	//"github.com/VineBalloon/nozobot/sounds"
	"github.com/VineBalloon/nozobot/voice"
	"github.com/bwmarrin/discordgo"
)

type Washi struct {
	Name string
}

func (w *Washi) Desc() string {
	return "Nozomi's washi washi will follow you into Voice as well!"
}

func (w *Washi) Roles() []string {
	return nil
}

func (w *Washi) Channels() []string {
	return nil
}

func (w *Washi) Handle(c *client.ClientState) error {
	s := c.Session
	m := c.Message

	_, err := s.ChannelMessageSend(m.ChannelID, "Washi Washi!")
	if err != nil {
		return err
	}

	// Create a voice room
	vr, err := voice.NewVoiceRoom(s, m.Message)
	if err != nil {
		return err
	}

	// Create a new sound map
	sm := map[string]*sounds.Sound{
		"1": sounds.NewSound("1", 10),
	}

	// Create a new sound collection with our sound map
	sounds.NewSoundCollection(w.Name, sm)

	// Connect to the voice room
	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Play a random sound
	vr.PlayRandom()

	// Close the voice connection
	vr.Close()
	return nil
}

func NewWashi(n string) *Washi {
	return &Washi{
		"Washi",
	}
}
