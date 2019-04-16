package handlers

import (
	"sort"
	"strings"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/utils"
)

// Help
// The command to generate help messages from other commands
type Help struct {
	name         string
	descriptions map[string]string
	prefix       string
}

func (h *Help) Name() string {
	return h.name
}

func (h *Help) Desc() string {
	return "Nozomi helps you write out this command!"
}

func (h *Help) Roles() []string {
	return nil
}

func (h *Help) Channels() []string {
	return nil
}

// MsgHandle
// Constructs the help message and cleanly formats it
func (h *Help) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	out := utils.Bold("Commands:\n")
	sorted := []string{}
	for name, desc := range h.descriptions {
		sorted = append(sorted, utils.Code(h.prefix+name)+": "+desc)
	}
	sort.Strings(sorted)
	out += strings.Join(sorted, "\n")

	_, err := s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		return err
	}
	return nil
}

// NewHelp
// Constructs a new help struct.
// Note: Requires the command prefix and router for dynamic help generation
func NewHelp(prefix string, r *map[string]Handler) *Help {
	// Register Descriptions from router
	desc := make(map[string]string)
	for cmd, handler := range *r {
		desc[cmd] = handler.Desc()
	}
	return &Help{
		"help",
		desc,
		prefix,
	}
}
