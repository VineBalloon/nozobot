package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/sounds"
	//"github.com/bwmarrin/discordgo"
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

func (w *Washi) Handle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	_, err := s.ChannelMessageSend(m.ChannelID, "Washi Washi!")
	if err != nil {
		return err
	}

	// Create a new sound map
	sm := map[string]*sounds.Sound{
		"1": sounds.NewSound("1", 10),
	}

	// Create a new sound collection with our sound map
	sc, err := sounds.NewSoundCollection(w.Name, sm)
	if err != nil {
		return err
	}

	// Create a voice room
	vr, err := client.NewVoiceRoom(s, m, sc)
	if err != nil {
		return err
	}

	// Connect to the voice room
	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Play a random sound
	vr.PlaySound()

	// Close the voice connection
	vr.Close()
	return nil
}

func NewWashi() *Washi {
	return &Washi{
		"Washi",
	}
}
