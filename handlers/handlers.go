package handlers

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Handler interface for commands to implement
type Handler interface {
	Desc() string
	Roles() []string
	Channels() []string
	Handle(*discordgo.Session, *discordgo.MessageCreate) error
}

type VoiceRoom struct {
	guild      string
	id         string
	Connection *discordgo.VoiceConnection
}

// Helper function to get guild and vc id from state and message
// Pass into constructor as `NewVoiceRoom(VoiceInfoFromMessage(session, user))
func VoiceInfoFromMessage(s *discordgo.Session, m *discordgo.Message) (string, string) {
	// Get the guild and guild ID
	mchannel, _ := s.Channel(m.ChannelID)
	guildId := mchannel.GuildID
	guild, _ := s.Guild(guildId)

	// Get channel id
	u := m.Author
	channel := ""
	for _, vs := range guild.VoiceStates {
		if vs.UserID == u.ID {
			channel = vs.ChannelID
		}
	}

	return guildId, channel
}

func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.guild, v.id, false, false)
	v.Connection = vc
	return err
}

func (v *VoiceRoom) Close() {
	v.Connection.Disconnect()
	v.Connection = nil
}

// Voice Room Constructor, always call with VoiceInfoFromMessage
func NewVoiceRoom(guild, channel string) (*VoiceRoom, error) {
	if channel == "" {
		return nil, errors.New("user not in voice channel")
	}

	return &VoiceRoom{
		guild:      guild,
		id:         channel,
		Connection: nil,
	}, nil
}
