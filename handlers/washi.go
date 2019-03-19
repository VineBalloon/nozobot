package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/sounds"
)

// TODO Washi
// The command to play 'washi washi' sounds in voice chat
// also maybe post a picture
type Washi struct {
	name        string
	description string
}

func (w *Washi) Name() string {
	return w.name
}

func (w *Washi) Desc() string {
	return w.description
}

func (w *Washi) Roles() []string {
	return nil
}

func (w *Washi) Channels() []string {
	return nil
}

// MsgHandle
// Tries to join voice, fails if already joined or if user isn't in one.
// TODO If no sound specified, plays a random sound from the washi sound collection.
// Like Junai, you can stop the sound with the StopSig channel.
// You probably won't though since most sounds are very short.
func (w *Washi) MsgHandle(cs *client.ClientState) error {
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
	sc, err := sounds.NewSoundCollection(w.name, sm)
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

	// Leave the voice connection
	vr.Leave()
	return nil
}

func NewWashi() *Washi {
	return &Washi{
		"Washi",
		"Nozomi's washi washi will follow you into Voice as well!",
	}
}
