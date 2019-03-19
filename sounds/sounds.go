package sounds

// Some stuff stolen from airhorn bot
// https://github.com/discordapp/airhornbot
// As well as example code from dca (check imports)

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jonas747/dca"
)

// SoundCollection aggregates Sound structs.
// N.B. Collections are hardcoded in each handler
type SoundCollection struct {
	Prefix string
	Sounds map[string]*Sound

	soundRange int
}

// CheckSound checks if the sound exists on disk
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

// EncodeName tries to encode a sound from the SoundCollection with that name
func (s *SoundCollection) EncodeName(name string) (*dca.EncodeSession, error) {
	// Get the named sound
	sound, found := s.Sounds[name]
	if !found {
		return nil, errors.New("encode: sound not found")
	}

	// Check for relative path things
	if strings.Contains(name, "..") {
		return nil, errors.New("encode: nice try")
	}

	// Check that the sound exists on disk
	err := s.CheckSound(name)
	if err != nil {
		return nil, err
	}

	// Encode the sound
	return sound.Encode(s.Prefix)
}

// EncodeRandom picks and encodes a random sound from the sound collection
func (s *SoundCollection) EncodeRandom() (*dca.EncodeSession, error) {
	i := 0
	number := randomRange(0, s.soundRange)

	for _, sound := range s.Sounds {
		i += sound.Weight

		if number < i {
			// Check that the sound exists on disk
			err := s.CheckSound(sound.Name)
			if err != nil {
				return nil, err
			}

			// return the encoder
			return sound.Encode(s.Prefix)
		}
	}
	return nil, errors.New("encoderandom: this should never happen")
}

// NewSoundCollection constructs a new SoundCollection.
// This should be called in each handler
func NewSoundCollection(prefix string, sounds map[string]*Sound) (*SoundCollection, error) {
	// Iterate through sounds and get the range for weights
	sr := 0
	for _, sound := range sounds {
		sr += sound.Weight
		if sound.Weight < 0 {
			return nil, errors.New("sounds: cannot parse negative weights")
		}
	}

	// Construct and return a new SoundCollection
	return &SoundCollection{
		Prefix:     prefix,
		Sounds:     sounds,
		soundRange: sr,
	}, nil
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
