package handlers

import (
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

func (j *Junai) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Attempt to join a voice room
	vr, err := NewVoiceRoom(VoiceInfoFromMessage(s, m.Message))
	if err != nil {
		return err
	}

	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Play junai lens
	sounds.JUNAI.Load()
	// TODO make this nicer
	sound := sounds.JUNAI.Random()
	sound.Play(vr.Connection, "junai")

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
