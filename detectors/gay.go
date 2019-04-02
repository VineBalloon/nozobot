package detectors

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
	/*
		whitelist = []string{
			"/r/anime",
			"/r/araragi",
			"/r/awwnime",
			"/r/akkordian",
			"/r/lovelive",
			"/r/washiwashi",
			"/r/wholesomeyuri",
		}
	*/
)

// Gay Detector for image messages
type Gay struct {
	name string
	desc string
}

// Name returns name of detector
func (g *Gay) Name() string {
	return g.name
}

// Desc returns description of detector
func (g *Gay) Desc() string {
	return g.desc
}

// MsgDetect detects messages to respond to
func (g *Gay) MsgDetect(cs *client.ClientState) error {
	s := cs.Session
	m := cs.Message
	request := strings.TrimRight(cs.Fullargs[0], " /")

	cl := &http.Client{}
	var imgurl string
	imgregex, _ := regexp.Compile("(http(s?):)([/|.|\\w|\\s|-])*\\.(?:jpg|gif|png)")
	switch {
	case imgregex.MatchString(request):
		imgurl = request

	case strings.HasPrefix(request, "https://www.reddit.com"):
		// Reddit links are case insensitive
		/* Don't use this yet
			request = strings.ToLower(request)
			for _, w := range whitelist {
				if strings.Index(request, w) != -1 {
					goto In
				}
			}
			return nil
		In:
		*/

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

	/* Add more cases here */
	default:
		return nil
	}

	// Just to be sure
	if !imgregex.MatchString(imgurl) {
		return nil
	}

	// Send typing, detect doesn't do this for us
	s.ChannelTyping(m.ChannelID)

	// Log so I can see how much of my API money we're using
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

	return nil
}

// Apiget Gets the API key required by the command
func (g *Gay) Apiget() error {
	var exists bool
	key, exists = os.LookupEnv("WAIFU")
	if !exists {
		return errors.New("Missing Waifu2x API Key: WAIFU")
	}
	return nil
}

// NewGay Constructs a new Gay structure
func NewGay() *Gay {
	return &Gay{
		"Gay",
		"Upscales images, also works with Reddit comment links for i.reddit",
	}
}
