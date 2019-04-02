package handlers

import (
	"strings"

	"github.com/VineBalloon/nozobot/client"
)

// Stat The command to change Nozomi's status
type Stat struct {
	name        string
	description string
}

// Name Returns name of the command
func (p *Stat) Name() string {
	return p.name
}

// Desc Returns description of the command
func (p *Stat) Desc() string {
	return p.description
}

// Roles Returns roles required by the command
func (p *Stat) Roles() []string {
	return nil
}

// Channels Returns channels required by the command
func (p *Stat) Channels() []string {
	return nil
}

// MsgHandle Changes Nozomi's status
func (p *Stat) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	args := cs.Args
	if len(args) == 0 {
		err := s.UpdateStatus(0, "")
		if err != nil {
			return err
		}
	}

	pre := strings.ToLower(args[0])
	stat := ""
	if pre == "listen" {
		if len(args) > 1 {
			stat = strings.Join(args[1:], " ")
		}
		err := s.UpdateListeningStatus(stat)
		if err != nil {
			return err
		}
		return nil
	}

	if pre == "play" {
		if len(args) > 1 {
			stat = strings.Join(args[1:], " ")
		}
	} else {
		if len(args) > 0 {
			stat = strings.Join(args, " ")
		}
	}
	err := s.UpdateStatus(0, stat)
	if err != nil {
		return err
	}
	return nil
}

// NewStat Constructor for Stat
func NewStat() *Stat {
	return &Stat{
		"Stat",
		"Changes Nozomi's Presence",
	}
}
