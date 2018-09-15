package sounds

// Stolen from airhorn bot
// https://github.com/discordapp/airhornbot

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// Collections are hardcoded, see bottom
type SoundCollection struct {
	Prefix string
	Sounds []*Sound

	soundRange int
}

func (sc *SoundCollection) Load() {
	for _, sound := range sc.Sounds {
		sc.soundRange += sound.Weight
	}
}

func (s *SoundCollection) Random() *Sound {
	var (
		i      int
		number int = randomRange(0, s.soundRange)
	)

	for _, sound := range s.Sounds {
		i += sound.Weight

		if number < i {
			return sound
		}
	}
	return nil
}

// Sound represents a sound clip
type Sound struct {
	Name   string
	Weight int
}

// Create a Sound struct
func NewSound(Name string, Weight int) *Sound {
	return &Sound{
		Name:   Name,
		Weight: Weight,
	}
}

// Plays this sound over the specified VoiceConnection
func (s *Sound) Play(vc *discordgo.VoiceConnection, cName string) error {
	// Send speaking packet over voice websocket
	err := vc.Speaking(true)
	if err != nil {
		return err
	}
	// Cleanup when we're done
	defer vc.Speaking(false)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	path := fmt.Sprintf("audio/%v_%v.wav",
		strings.ToLower(cName), strings.ToLower(s.Name))

	fmt.Println("Playing:", path)

	enc, err := dca.EncodeFile(path, opts)
	if err != nil {
		return err
	}

	done := make(chan error)
	_ = dca.NewStream(enc, vc, done)
	ticker := time.NewTicker(time.Second)

	for {

		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				return err
			}

			// Clean up
			enc.Truncate()
			return nil
		case <-ticker.C:
			//stats := enc.Stats()
			//playPos := stream.PlaybackPosition()

		}
	}
}

// Returns a random integer between min and max
func randomRange(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// Hardcoded sound collections
var JUNAI *SoundCollection = &SoundCollection{
	Prefix: "junai",
	Sounds: []*Sound{
		NewSound("lens", 100),
	},
}
