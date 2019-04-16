package handlers

import (
	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/sounds"
	"github.com/VineBalloon/nozobot/utils"
)

// Junai
// The command to play Junai Lens in voice chat
type Junai struct {
	name        string
	description string
}

func (j *Junai) Name() string {
	return j.name
}

func (j *Junai) Desc() string {
	return j.description
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
func (j *Junai) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	// Create a new sound map
	sm := map[string]*sounds.Sound{
		"lens": sounds.NewSound("lens", 100),
	}

	// Create a new sound collection with our sound map
	sc, err := sounds.NewSoundCollection(j.name, sm)
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
	_, err = s.ChannelMessageSend(m.ChannelID, utils.Bold("Ikuyoooo!"))
	if err != nil {
		return err
	}

	// Play junai lens
	cs.Voice.PlaySound("lens")
	return nil
}

func NewJunai() *Junai {
	return &Junai{
		"junai",
		"Nozomi sings Junai Lens for you!",
	}
}
