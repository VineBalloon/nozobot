package handlers

import (
	"time"

	//"github.com/VineBalloon/nozobot/helpers"
	//"github.com/VineBalloon/nozobot/sounds"
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

func (w *Washi) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Washi Washi!")
	if err != nil {
		return err
	}

	// Attempt to join a voice room
	vr, err := NewVoiceRoom(VoiceInfoFromMessage(s, m.Message))
	if err != nil {
		return err
	}

	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Sleep for 5 seconds
	// TODO make nozomi play the audio
	time.Sleep(time.Second * 5)

	// Close the voice connection
	vr.Close()
	return nil
}

func NewWashi(n string) *Washi {
	return &Washi{
		"Washi",
	}
}
