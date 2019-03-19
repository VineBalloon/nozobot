package client

// This file contains methods related to transmitting over the voice connection

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jonas747/dca"
)

// PlaySound
// Plays a named sound from the sound collection
// and a random sound if no name is given
func (v *VoiceRoom) PlaySound(names ...string) error {
	// Check if we are already playing something
	if v.Stream != nil {
		return errors.New("playsound: already playing something")
	}

	// Check if voice room has a sound collection
	if v.Sounds == nil {
		return errors.New("playsound: no sound collection in voice room")
	}

	// Check our slice of names
	var enc *dca.EncodeSession
	var err error
	switch {
	case len(names) == 0:
		// Get random sound
		enc, err = v.Sounds.EncodeRandom()
		if err != nil {
			return err
		}

	case len(names) == 1:
		// Get the named sound
		enc, err = v.Sounds.EncodeName(names[0])
		if err != nil {
			return err
		}

	default:
		// Only handle 0 or 1 sounds
		return errors.New("playsound: invalid number of sounds")
	}

	// Send speaking packet over voice websocket
	err = v.Connection.Speaking(true)
	if err != nil {
		return err
	}

	// Start a new stream from the encoding session
	// to the discord voice connection
	done := make(chan error)
	v.Stream = dca.NewStream(enc, v.Connection, done)

	// Make signal channel to stop
	v.StopSig = make(chan struct{})

	// Some shit from dca example iunno probably use it for presence
	ticker := time.NewTicker(time.Second)

	// Async audio loop
	for {
		select {
		case <-v.StopSig: // Received stop from voice
			v.StopSig = nil
			fmt.Println("Received Stop")

			// Stream stopped, clean up but don't dc
			fmt.Println("Cleaning up")
			enc.Cleanup()
			v.Stream = nil

			return nil

		case err := <-done: // Received done from encoder
			if err != nil && err != io.EOF {
				// Not 'done', some other error
				v.Stream = nil
				v.StopSig = nil
				enc.Cleanup()
				return err
			}
			fmt.Println("Received Done")

			// Stream done, clean up but don't dc
			fmt.Println("Cleaning up")
			enc.Cleanup()
			v.Stream = nil

			return nil

		case <-ticker.C: // Ticker fires off every second to update stats
			//stats := enc.Stats()
			//playPos := v.Stream.PlaybackPosition()
		}
	}

	return nil
}

// Pause
// Attempts to pause the current stream.
func (v *VoiceRoom) Pause() error {
	s := v.Stream
	if s == nil {
		return errors.New("pause: no voice stream to pause")
	}

	if s.Paused() {
		return errors.New("pause: already paused")
	}

	v.Stream.SetPaused(true)

	return nil
}

// UnPause
// Attempts to unpause the current stream.
func (v *VoiceRoom) UnPause() error {
	s := v.Stream
	if s == nil {
		return errors.New("unpause: no voice stream to pause")
	}

	if !s.Paused() {
		return errors.New("unpause: already unpaused")
	}

	v.Stream.SetPaused(false)

	return nil
}

// Stop
// Gracefully end the streaming session without disconnecting.
func (v *VoiceRoom) Stop() error {
	// Check channel
	if v.StopSig == nil {
		return errors.New("stop: no voice stream to stop")
	}

	// Send stop signal to the StopSig channel
	close(v.StopSig)
	return nil
}
