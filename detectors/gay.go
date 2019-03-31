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
)

// Gay detector for image messages
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
		request += ".json"
		//req, err := http.NewRequest("GET", "http://httpbin.org/user-agent", nil)
		req, err := http.NewRequest("GET", request, nil)
		if err != nil {
			log.Println(err)
			return nil
		}

		req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")
		resp, err := cl.Do(req)
		if err != nil {
			log.Println(err)
			return nil
		}

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil
		}

		val, _, _, err := jsonparser.Get( // Magic
			b,
			"[0]", "data", "children",
			"[0]", "data", "url",
		)
		if err != nil {
			log.Println(err)
			return nil
		}

		imgurl = string(val)
		if !imgregex.MatchString(imgurl) {
			return nil
		}

	/* Add more cases here */
	default:
		return nil
	}

	// Send typing, detect doesn't do this for us
	s.ChannelTyping(m.ChannelID)

	log.Println("Requesting url: ", imgurl)
	form := url.Values{}
	form.Add("image", imgurl)
	req, err := http.NewRequest("POST", waifu, strings.NewReader(form.Encode()))
	if err != nil {
		log.Println(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("api-key", key)
	resp, err := cl.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	val, _, _, err := jsonparser.Get(b, "output_url")
	out := string(val)
	_, err = s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		log.Println(err)
		return nil
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
