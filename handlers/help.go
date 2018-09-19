package handlers

import (
	"sort"
	"strings"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/helpers"
	//"github.com/bwmarrin/discordgo"
)

// Help
// The command to generate help messages from other commands
type Help struct {
	Name         string
	descriptions map[string]string
	prefix       string
}

// AddDesc
// Generates descriptions given the router's string->handler map.
// This method is unique to the Help command
func (h *Help) AddDesc(r *map[string]Handler) {
	h.descriptions = make(map[string]string)
	for cmd, handler := range *r {
		h.descriptions[cmd] = handler.Desc()
	}
}

// Handle
// Constructs the help message and cleanly formats it
func (h *Help) Handle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message

	out := helpers.Bold("Commands:\n")
	sorted := []string{}
	for name, desc := range h.descriptions {
		sorted = append(sorted, helpers.Code(h.prefix+name)+": "+desc)
	}
	sort.Strings(sorted)
	out += strings.Join(sorted, "\n")

	_, err := s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		return err
	}
	return nil
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

// NewHelp
// Constructs a new help struct. Requires the command prefix.
func NewHelp(prefix string) *Help {
	return &Help{
		"help",
		nil,
		prefix,
	}
}
