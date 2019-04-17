package detectors

import (
	"math/rand"
	"strings"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/utils"
)

var (
	chotto = []string{
		"Wait, ",
		"Dame dame, ",
		"Chotto, ",
		"Chotto matte kudosai, ",
		utils.Italics("taps shoulder") + " nee nee~, ",
	}
	illegal = []string{
		"loli",
		"drug",
		"hentai",
		"rule 34",
	}
)

// Illegal Detector for illegal messages
type Illegal struct {
	name string
	desc string
}

// Name Returns name of detector
func (i *Illegal) Name() string {
	return i.name
}

// Desc Returns description of detector
func (i *Illegal) Desc() string {
	return i.desc
}

// MsgDetect Detects messages to respond to
func (i *Illegal) MsgDetect(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	lowered := strings.ToLower(cs.Message.Content)
	if utils.Slicehassubstring(lowered, illegal) {
		out := chotto[rand.Int()%len(chotto)] + "that's illegal!"
		s.ChannelMessageSend(m.ChannelID, out)
		return nil
	}
	return nil
}

// Apiget Gets the API key required by the command
func (i *Illegal) Apiget() error {
	return nil
}

// NewIllegal Constructs a new Illegal structure
func NewIllegal() *Illegal {
	return &Illegal{
		"Illegal",
		"Nozomi detects illegal things 👮‍♀️",
	}
}