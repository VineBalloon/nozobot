package client

// This file contains the wrapper structs for our discord client's
// session and voice connection, and methods to manipulate the Session

import "github.com/bwmarrin/discordgo"

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
