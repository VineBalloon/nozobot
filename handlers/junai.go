package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/helpers"
	"github.com/VineBalloon/nozobot/sounds"
)

// Junai
// The command to play Junai Lens in voice chat
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

// Handle
// Tries to join voice, fails if already joined or if user isn't in one.
// Plays Junai Lens on successful join.
// NB: This function keeps an async loop alive until the stream ends.
// You can stop it completely using the StopSig channel
func (j *Junai) Handle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	// Create a new sound map
	sm := map[string]*sounds.Sound{
		"lens": sounds.NewSound("lens", 100),
	}

	// Create a new sound collection with our sound map
	sc, err := sounds.NewSoundCollection(j.Name, sm)
	if err != nil {
		return err
	}

	// Create a new voice room with the SoundCollection
	cs.Voice, err = client.NewVoiceRoom(s, m, sc)
	if err != nil {
		return err
	}

	// Connect to the voice channel
	err = cs.Voice.Connect(s)
	if err != nil {
		return err
	}

	// Signal to the people that we are about to get rowdy
	_, err = s.ChannelMessageSend(m.ChannelID, helpers.Bold("Ikuyoooo!"))
	if err != nil {
		cs.Voice.Close()
		return err
	}

	// Play junai lens
	cs.Voice.PlaySound("lens")

	// Close the voice connection when we're done
	cs.Voice.Close()
	return nil
}

func NewJunai() *Junai {
	return &Junai{
		"Junai",
	}
}
