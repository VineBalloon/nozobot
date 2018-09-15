package handlers

import (
	"sort"
	"strings"

	"github.com/VineBalloon/nozobot/helpers"
	"github.com/bwmarrin/discordgo"
)

type Help struct {
	Name         string
	descriptions map[string]string
	prefix       string
}

func (h *Help) AddDesc(r *map[string]Handler) {
	h.descriptions = make(map[string]string)
	for cmd, handler := range *r {
		h.descriptions[cmd] = handler.Desc()
	}
}

func (h *Help) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

func NewHelp(n string, p string) *Help {
	return &Help{
		"help",
		nil,
		p,
	}
}
