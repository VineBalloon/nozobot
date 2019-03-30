package handlers

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/VineBalloon/nozobot/client"
	"github.com/buger/jsonparser"
)

var (
	waifu = "https://api.deepai.org/api/waifu2x"
	key   string
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
	return []string{"gay", "boi"}
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

	cl := &http.Client{}
	for r := range requests {
		request := strings.TrimRight(requests[r], " /")
		imgreg, err := regexp.Compile("(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")
		if err != nil {
			return err
		}

		var imgurl string
		switch {
		case imgreg.MatchString(request):
			imgurl = request

		case strings.HasPrefix(request, "https://www.reddit.com"):
			request += ".json"
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

			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			val, _, _, err := jsonparser.Get( // Magic
				b,
				"[0]", "data", "children",
				"[0]", "data", "url",
			)
			if err != nil {
				return err
			}

			imgurl = string(val)

		// Add more cases here
		default:
			return errors.New("gay: not an image or reddit url")
		}

		log.Println("Requesting url: ", imgurl)
		form := url.Values{}
		form.Add("image", imgurl)
		req, err := http.NewRequest("POST", waifu, strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("api-key", key)
		resp, err := cl.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		val, _, _, err := jsonparser.Get(b, "output_url")
		out := string(val)
		_, err = s.ChannelMessageSend(m.ChannelID, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGay() *Gay {
	return &Gay{
		"Gay",
		"Upscales images, also works with Reddit comment links for i.reddit",
	}
}

func init() {
	var exists bool
	key, exists = os.LookupEnv("WAIFU")
	if !exists {
		log.Fatal("Missing Waifu2x API Key: WAIFU")
	}
}
