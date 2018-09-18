package client

import (
	"errors"

	"github.com/VineBalloon/nozobot/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// ClientState is a wrapper around the discordgo session,
// the VoiceRoom wrapper, and the last message received
type ClientState struct {
	Session *discordgo.Session
	Voice   *VoiceRoom
	Message *discordgo.Message
}

// UpdateMessage updates the discord message of ClientStatestruct
func (c *ClientState) UpdateMessage(m *discordgo.Message) {
	c.Message = m
}

// AddVoice adds a VoiceRoom to ClientState
func (c *ClientState) AddVoice(vr *VoiceRoom) {
	c.Voice = vr
}

// NewClientState constructs a new ClientState
func NewClientState(s *discordgo.Session, m *discordgo.Message) *ClientState {
	return &ClientState{
		Session: s,
		Voice:   nil,
		Message: m,
	}
}

// VoiceRoom is a wrapper around a voice connection
// and a sound collection, with methods to manipulate both
type VoiceRoom struct {
	guild      string
	id         string
	Connection *discordgo.VoiceConnection
	Stream     *dca.StreamingSession
	Sounds     *sounds.SoundCollection
}

// Connect connects the session client to the voice room
func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.guild, v.id, false, false)
	v.Connection = vc
	return err
}

// Play plays the encoding session over the voice channel
func (v *VoiceRoom) Play(enc *dca.EncodeSession) {
	// Send speaking packet over voice websocket
	err := vc.Speaking(true)
	if err != nil {
		return err
	}
	// Cleanup when we're done
	defer vc.Speaking(false)

	// Start a new stream from the encoding session
	// to the discord voice connection
	done := make(chan error)
	v.Stream = dca.NewStream(enc, vc, done)
	ticker := time.NewTicker(time.Second)

	// Async event loop for our audio stream
	for {
		select {
		case err := <-done:
			// done channel has been sent an error, handle it
			if err != nil && err != io.EOF {
				// Not 'done', some other error
				return err
			}

			// Stream done, clean up
			enc.Cleanup()
			return nil
		case <-ticker.C:
			// Ticker when not done
			//stats := enc.Stats()
			playPos := v.Stream.PlaybackPosition()
		}
	}
}

// PlayRandom plays a random sound under the command's prefix
func (v *VoiceRoom) PlayRandom() {
	v.Sounds.PlayRandom(v.Connection)
}

// Close closes the voice connection with discord and flushes the connection pointer
func (v *VoiceRoom) Close() {
	v.Connection.Disconnect()
	v.Connection = nil
}

// NewVoiceRoom creates a new Voice Room using the current session and message received.
// It uses the helper function VoiceInfoFromMessage
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

// VoiceInfoFromMessage is a helper function to get guild and vc id from the current discord session
// and message received
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
		return "", "", errors.New("user not in voice channel")
	}

	return guildId, channel, nil
}
