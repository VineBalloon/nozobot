package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/VineBalloon/nozobot/client"
)

// Gay
// The command to get the image from a reddit link
type Gay struct {
	name        string
	description string
}

func (g *Gay) Name() string {
	return g.name
}

func (g *Gay) Desc() string {
	return g.description
}

func (g *Gay) Roles() []string {
	return nil
}

func (g *Gay) Channels() []string {
	return nil
}

// MsgHandle
// Gets i.reddit image from reddit comments link
func (g *Gay) MsgHandle(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	requests := cs.Arguments
	if len(requests) == 0 {
		return errors.New("argparse: not enough arguments")
	}

	for r := range requests {
		request := requests[r]
		if !strings.HasPrefix(request, "https://www.reddit.com") {
			return errors.New("gay: not a valid reddit url!")
		}
		strings.TrimRight(request, " /")
		resp, err := http.Get(request + ".json")
		if err != nil {
			return err
		}
		// TODO
		url := "Placeholder"
		fmt.Println(resp)

		var j interface{}
		err = json.Unmarshal(resp, &j)
		if err != nil {
			return err
		}
		fmt.Printf("%+v", j)

		_, err = s.ChannelMessageSend(m.ChannelID, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGay() *Gay {
	return &Gay{
		"Gay",
		"Gets i.reddit image from reddit comments link",
	}
}
