package client

// This file contains the wrapper structs for our discord client's
// session and voice connection, and methods to manipulate the Session

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ClientState
// A wrapper around the discordgo session,
// the VoiceRoom wrapper, and the last message received
type ClientState struct {
	Session   *discordgo.Session /* The current session */
	Voice     *VoiceRoom         /* Voice wrapper */
	Message   *discordgo.Message /* Last message received */
	Arguments []string           /* Arguments from the message */
}

// UpdateState
// Updates the current discord session and last message received.
func (c *ClientState) UpdateState(s *discordgo.Session, m *discordgo.Message) {
	c.Session = s
	c.Message = m
	c.Arguments = strings.Split(m.Content, " ")[1:]
}

// AddVoice
// Adds a VoiceRoom to ClientState
func (c *ClientState) AddVoice(vr *VoiceRoom) {
	c.Voice = vr
}

// Close
// Closes all connections if there are any
func (c *ClientState) Close() error {
	// Leave voice connection, ignore errors
	if c.Voice == nil {
		return errors.New("close: no voice room to close")
	}
	c.Voice.Leave()

	// Close discord connection
	err := c.Session.Close()
	if err != nil {
		return err
	}
	return nil
}

// NewClientState constructs a new ClientState
func NewClientState() *ClientState {
	return &ClientState{
		nil,
		nil,
		nil,
		nil,
	}
}
