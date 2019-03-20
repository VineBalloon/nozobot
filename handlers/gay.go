package handlers

import (
	//"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/VineBalloon/nozobot/client"
	"github.com/buger/jsonparser"
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
	fmt.Println(requests)

	cl := &http.Client{}
	for r := range requests {
		request := strings.TrimRight(requests[r], " /") + ".json"
		fmt.Println(request)

		if !strings.HasPrefix(request, "https://www.reddit.com") {
			return errors.New("gay: not a valid reddit url!")
		}

		//req, err := http.NewRequest("GET", "http://httpbin.org/user-agent", nil)
		req, err := http.NewRequest("GET", request, nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

		resp, err := cl.Do(req)
		if err != nil {
			return err
		}
		fmt.Println(resp)

		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		//var j []interface{}
		//err = json.NewDecoder(resp.Body).Decode(&j)
		if err != nil {
			return err
		}
		/*
			url := "Placeholder"
			out, err := json.MarshalIndent(j, "", "    ")
			if err != nil {
				return err
			}
		*/
		// slice of map of string to map of string to slice of map of string to map of string to string
		//data := j[0].(map[string]map[string][]map[string]map[string]string)
		//jdata := j[0].(map[string]interface{}).(map[string]interface{})
		//jurl := data["data"]["children"][0]["data"]["url"]
		val, _, _, err := jsonparser.Get(bytes, "[0]", "data", "children", "[0]", "data", "url")
		if err != nil {
			return err
		}
		url := string(val)
		fmt.Println(url)

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
