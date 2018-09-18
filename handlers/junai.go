package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/helpers"
	"github.com/VineBalloon/nozobot/sounds"
	"github.com/bwmarrin/discordgo"
)

type Junai struct {
	Name string
}

func (j *Junai) Desc() string {
	return "Nozomi sings Junai Lens for you!"
}

func (j *Junai) Roles() []string {
	return nil
}

func (j *Junai) Channels() []string {
	return nil
}

func (j *Junai) Handle(c *client.ClientState) error {
	s := c.Session
	m := c.Message

	// Create a new voice room
	vr, err := client.NewVoiceRoom(s, m, sc)
	if err != nil {
		return err
	}

	// Create a new sound map
	sm := map[string]*sounds.Sound{
		"lens": sounds.NewSound("lens", 100),
	}

	// Create a new sound collection with our sound map
	sc := sounds.NewSoundCollection(j.Name, sm)

	// Connect to the voice channel
	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Play junai lens
	vr.PlayName("lens")

	// Signal to the people that we are about to get rowdy
	_, err = s.ChannelMessageSend(m.ChannelID, helpers.Bold("Ikuyoooo!"))
	if err != nil {
		return err
	}

	// Close the voice connection
	vr.Close()
	return nil
}

func NewJunai(n string) *Junai {
	return &Junai{
		"Junai",
	}
}
