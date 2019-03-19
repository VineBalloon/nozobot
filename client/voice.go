package client

// This file contains the wrapper struct for voice connections
// and methods to manipulate the VoiceConnection

import (
	"errors"

	"github.com/VineBalloon/nozobot/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// VoiceRoom
// A wrapper around a voice connection and a sound collection,
// with methods to manipulate both.
type VoiceRoom struct {
	Guild      string
	Id         string
	Connection *discordgo.VoiceConnection

	Sounds  *sounds.SoundCollection
	Stream  *dca.StreamingSession
	StopSig chan struct{}
}

// Connect
// Connects the session client to the voice room.
func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.Guild, v.Id, false, false)
	v.Connection = vc
	return err
}

// Leave
// Leaves the voice connection with discord and flushes the connection pointer.
func (v *VoiceRoom) Leave() error {
	// Check connection
	if v.Connection == nil {
		return errors.New("leave: nothing to leave")
	}

	// Stop
	v.Stop()

	// Close voice connection
	v.Connection.Close()

	// Reset connection pointer
	v.Connection = nil
	return nil
}

// NewVoiceRoom
// Creates a new Voice Room using the current session and message received.
func NewVoiceRoom(s *discordgo.Session, m *discordgo.Message, sounds *sounds.SoundCollection) (*VoiceRoom, error) {

	// Get the voice info from the message
	guild, channel, err := VoiceInfoFromMessage(s, m)
	if err != nil {
		return nil, err
	}

	return &VoiceRoom{
		Guild:      guild,
		Id:         channel,
		Connection: nil,

		Sounds:  sounds,
		Stream:  nil,
		StopSig: nil,
	}, nil
}

// VoiceInfoFromMessage
// A helper function to get guild and vc id
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
