package sounds

// Some stuff stolen from airhorn bot
// https://github.com/discordapp/airhornbot
// As well as example code from dca

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// SoundCollection aggregates Sound structs.
// N.B. Collections are hardcoded in each handler
type SoundCollection struct {
	Prefix string
	Sounds map[string]*Sound

	soundRange int
}

// CheckSound checks if the sound exists in the expected path
func (s *SoundCollection) CheckSound(name string) error {
	// Get path
	path := fmt.Sprintf("audio/%v/%v.wav",
		strings.ToLower(s.Prefix), strings.ToLower(name))

	// Confirm that path exists
	_, err := os.Open(path)
	if err != nil {
		return err
	}
	return nil
}

// Play a random sound from this collection
func (s *SoundCollection) PlayRandom(vc *discordgo.VoiceConnection) *Sound {
	var (
		i      int
		number int = randomRange(0, s.soundRange)
	)

	for _, sound := range s.Sounds {
		i += sound.Weight

		if number < i {
			//sound.Play(
		}
	}
	return nil
}

// EncodeName tries to encode a sound from the SoundCollection with that name
func (s *SoundCollection) EncodeName(name string) (*dca.EncodeSession, error) {
	// Get the named sound
	sound, found := s.Sounds[name]
	if !found {
		return nil, errors.New("sound not found")
	}

	// Check that the sound exists on disk
	err := s.CheckSound(name)
	if err != nil {
		return nil, err
	}

	// Encode the sound
	return sound.Encode(s.Prefix)
}

// NewSoundCollection constructs a new SoundCollection.
// This should be called in each handler
func NewSoundCollection(prefix string, sounds map[string]*sounds.Sound) *sounds.SoundCollection {
	// Iterate through sounds and get the range for weights
	sr := 0
	for _, sound := range sounds {
		sr += sound.Weight
	}

	// Construct and return a new SoundCollection
	return &sounds.SoundCollection{
		Prefix:     prefix,
		Sounds:     sounds,
		soundRange: sr,
	}
}

// Sound represents a sound clip
type Sound struct {
	Name   string
	Weight int
}

// Encode encodes the sound on disk and returns the encoding session
func (s *Sound) Encode(prefix string) (*dca.EncodeSession, error) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	// Get path for audio file
	path := fmt.Sprintf("audio/%v/%v.wav",
		strings.ToLower(prefix), strings.ToLower(s.Name))

	//fmt.Println("Playing:", path)

	// Create and return a new encoder
	// the encoder encodes audio file to an Opusreader stream
	return dca.EncodeFile(path, opts)
}

// NewSound creates a Sound struct
// weight represents how likely the sound will be chosen at random
func NewSound(Name string, Weight int) *Sound {
	return &Sound{
		Name:   Name,
		Weight: Weight,
	}
}

// randomRange returns a random integer between min and max
func randomRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
