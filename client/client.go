package client

// This file contains the wrapper structs for our discord client's
// session and voice connection, and methods to manipulate the Session

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ClientState is a wrapper around the discordgo session,
// the VoiceRoom wrapper, and the last message received
type ClientState struct {
	Session   *discordgo.Session
	Voice     *VoiceRoom
	Message   *discordgo.Message
	Arguments []string
}

// UpdateSession
// Updates the current discord session and message.
// It also parses the space separated arguments
func (c *ClientState) UpdateSession(s *discordgo.Session, m *discordgo.Message) {
	c.Session = s
	c.Message = m
	c.Arguments = strings.Split(m.Content, " ")[1:]
}

// AddVoice
// Adds a VoiceRoom to ClientState
func (c *ClientState) AddVoice(vr *VoiceRoom) {
	c.Voice = vr
}

// StopStream
// Stops the audio stream if it exists.
// Does not disconnect the bot.
func (c *ClientState) StopStream() error {
	if c.Voice == nil {
		return errors.New("client: no voice room")
	}
	err := c.Voice.Stop()
	if err != nil {
		return err
	}
	return nil
}

// Close
// Close all connections if there are any
func (c *ClientState) Close() {
	c.StopStream()
	c.Session.Close()
	c.Voice.Close()
}

// NewClientState constructs a new ClientState
func NewClientState(s *discordgo.Session, m *discordgo.Message) *ClientState {
	return &ClientState{
		Session: s,
		Voice:   nil,
		Message: m,
	}
}
