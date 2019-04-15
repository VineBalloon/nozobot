package helpers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Italics Surrounds with italics escapes
func Italics(s string) string {
	return "*" + s + "*"
}

// Bold Surrounds with bold escapes
func Bold(s string) string {
	return "**" + s + "**"
}

// Code Surrounds with code escapes
func Code(s string) string {
	return "`" + s + "`"
}

// Spoiler Surrounds with spoiler escapes
func Spoiler(s string) string {
	return "||" + s + "||"
}

// Noembed Surrounds with no-embed escapes
func Noembed(s string) string {
	return "<" + s + ">"
}

// At Surrounds ID with @ escapes
func At(s string) string {
	return "<@" + s + ">"
}

// Chan Surrounds ID with chan escapes
func Chan(s string) string {
	return "<#" + s + ">"
}

// Stringinslice Checks if string is in the slice
func Stringinslice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// Slicehassubstring Checks if any element of the slice is a substring of str
func Slicehassubstring(str string, slice []string) bool {
	for _, s := range slice {
		if strings.Index(str, s) != -1 {
			return true
		}
	}
	return false
}

// Hasroles Checks if the author has the required roles
func Hasroles(s *discordgo.Session, m *discordgo.Message, roles []string) (bool, error) {
	if len(roles) == 0 {
		return true, nil
	}
	// Get member
	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err != nil {
		member, err = s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			return false, err
		}
	}
	groles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		return false, err
	}
	rolesrequired := []string{}
	for _, r := range roles {
		for _, gr := range groles {
			if strings.ToLower(gr.Name) == r {
				rolesrequired = append(rolesrequired, gr.ID)
			}
		}
	}
	mroles := member.Roles
	for _, rr := range rolesrequired {
		if Stringinslice(rr, mroles) {
			return true, nil
		}
	}
	return false, nil
}

// Inchannel Checks if message was sent in the required channels
func Inchannel(s *discordgo.Session, m *discordgo.Message, channels []string) (bool, error) {
	if len(channels) == 0 {
		return true, nil
	}
	curr, err := s.Channel(m.ChannelID)
	if err != nil {
		return false, err
	}
	for _, c := range channels {
		if c == curr.Name {
			return true, nil
		}
	}
	return false, nil
}
