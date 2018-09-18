package client

// This file contains the wrapper struct for voice connections
// and methods to manipulate the VoiceConnection
// The voice transmission is tightly coupled with the sounds package

import (
	"errors"
	"io"

	"github.com/VineBalloon/nozobot/sounds"
	"github.com/jonas747/dca"
)

// STOP is a global stop error for our async audio loop.
var STOP error = errors.New("stop")

// VoiceRoom is a wrapper around a voice connection
// and a sound collection, with methods to manipulate both.
type VoiceRoom struct {
	guild      string
	id         string
	Connection *discordgo.VoiceConnection

	Sounds *sounds.SoundCollection
	Stream *dca.StreamingSession
	Done   chan error
}

// Connect connects the session client to the voice room.
func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.guild, v.id, false, false)
	v.Connection = vc
	return err
}

// PlaySound plays a random sound from the sound collection if no name is given
// and the named sound if one was given.
func (v *VoiceRoom) PlaySound(names ...string) error {
	// Check if we are already playing something
	if v.Stream != nil {
		return errors.New("client: already playing something")
	}

	// Check if voice room has a sound collection
	if v.Sounds == nil {
		return errors.New("client: no sound collection in voice room")
	}

	// Check our slice of names
	switch {
	case len(names) == 0:
		// Get random sound
		enc, err := v.Sounds.EncodeRandom()
	case len(names) == 1:
		// Get the named sound
		enc, err := v.Sounds.EncodeName(name[0])
	default:
		// Only handle 0 or 1 sounds
		return errors.New("playsound: invalid number of sounds")
	}

	if err != nil {
		return err
	}

	// Send speaking packet over voice websocket
	err := v.Connection.Speaking(true)
	if err != nil {
		return err
	}

	// Cleanup when we're done
	defer v.Connection.Speaking(false)

	// Start a new stream from the encoding session
	// to the discord voice connection
	v.Done = make(chan error)
	v.Stream = dca.NewStream(enc, vc, v.Done)
	ticker := time.NewTicker(time.Second)

	// Async audio loop
	for {
		select {
		case err := <-v.Done:
			// done channel has been sent an error, handle it
			if err != nil && err != io.EOF && err != STOP {
				// Not 'done', some other error
				v.Stream = nil
				v.Done = nil
				return err
			}

			// Stream done, clean up encoder
			enc.Cleanup()
		case <-ticker.C:
			// Ticker when not done
			//stats := enc.Stats()
			playPos := v.Stream.PlaybackPosition()
		}
	}

	// Cleanup the stream
	v.Stream = nil
	return nil
}

// Pause attempts to pause the current stream.
func (v *VoiceRoom) Pause() {
	s := v.Stream
	if s == nil {
		return errors.New("pause: no voice stream to pause")
	}

	if s.Paused() {
		return erros.New("pause: already paused")
	}

	v.Stream.SetPaused(true)

	return nil
}

// UnPause attempts to pause the current stream.
func (v *VoiceRoom) UnPause() error {
	s := v.Stream
	if s == nil {
		return errors.New("pause: no voice stream to pause")
	}

	if !s.Paused() {
		return erros.New("unpause: already unpaused")
	}

	v.Stream.SetPaused(false)

	return nil
}

// Stop tries to gracefully end the streaming session without disconnecting.
func (v *VoiceRoom) Stop() error {
	if v.Stream == nil {
		return errors.New("stop: no voice stream to stop")
	}
	// Send stop signal to the Done channel
	v.Done <- STOP
	return nil
}

// Close closes the voice connection with discord and flushes the connection pointer.
func (v *VoiceRoom) Close() error {
	err := v.Connection.Disconnect()
	v.Connection = nil
	return err
}

// NewVoiceRoom creates a new Voice Room using the current session and message received.
// It uses the helper function VoiceInfoFromMessage. Use Connect to connect to the channel.
func NewVoiceRoom(s *discordgo.Session, m *discordgo.Message, sounds *sounds.SoundCollection) (*VoiceRoom, error) {

	// Get the voice info from the message
	guild, channel, err := VoiceInfoFromMessage(s, m)
	if err != nil {
		return nil, err
	}

	return &VoiceRoom{
		guild:      guild,
		id:         channel,
		Connection: nil,
		Sounds:     sounds,
	}, nil
}

// VoiceInfoFromMessage is a helper function to get guild and vc id
// from the current discord session and message received
func VoiceInfoFromMessage(s *discordgo.Session, m *discordgo.Message) (string, string, error) {
	// Get the guild and guild ID
	mchannel, err := s.Channel(m.ChannelID)
	if err != nil {
		return "", "", err
	}

	guildId := mchannel.GuildID

	guild, err := s.Guild(guildId)
	if err != nil {
		return "", "", err
	}

	// Get channel id
	u := m.Author
	channel := ""
	for _, vs := range guild.VoiceStates {
		if vs.UserID == u.ID {
			channel = vs.ChannelID
		}
	}

	// Throw error when user isn't in a voice channel
	if channel == "" {
		return "", "", errors.New("voiceinfo: user not in voice channel")
	}

	return guildId, channel, nil
}
