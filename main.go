// Main of Nozobot
// Modified from github.com/bwmarrin/disgord

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/handlers"
	h "github.com/VineBalloon/nozobot/helpers"
	"github.com/bwmarrin/discordgo"
)

// Global vars
var (
	token  string
	prefix = "!"
)

// Router
// Struct to hold our string->router mappings
type Router struct {
	routes map[string]handlers.Handler
}

// AddHandler
// Adds a string->handler mapping to the router
func (r *Router) AddHandler(cmd string, handler handlers.Handler) {
	r.routes[strings.ToLower(cmd)] = handler
}

// Route
// Returns the handler
func (r *Router) Route(cmd string) (handlers.Handler, bool) {
	routes := r.routes[cmd]
	if routes == nil {
		return nil, false
	}
	return routes, true
}

// Run
// Runs the discordgo handler and routes to our handlers
func (r *Router) Run(d *discordgo.Session) {
	cs := client.NewClientState(nil, nil)
	// Add anonymous function to route messages to handlers
	// this gets called when a message is read by the bot
	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Don't talk to yourself
		if m.Author.ID == s.State.User.ID {
			return
		}

		message := strings.TrimSpace(m.Content)
		if !strings.HasPrefix(message, prefix) {
			return
		}

		// Parse message
		command := strings.ToLower(strings.TrimLeft(message, prefix))

		//fmt.Printf("Received message {%s}\n", message)

		// Split the command into 2 substrings
		args := strings.SplitN(command, " ", 2)

		// Find handler using router
		handler, found := r.Route(args[0])
		if !found {
			s.ChannelMessageSend(m.ChannelID,
				h.Italics("Ara Ara: Unknown command, use "+h.Code(prefix+"help")))
			return
		}

		// Update the ClientState
		cs.UpdateSession(s, m.Message)

		// Call handler method
		err := handler.Handle(cs)

		// Got error from handler, throw it
		if err != nil {
			fmt.Println(err.Error())
			s.ChannelMessageSend(m.ChannelID,
				h.Italics("Ara Ara:"+strings.Split(err.Error(), ":")[1:]))
			return
		}
	})

	// Don't close the connection, wait for a kill signal
	fmt.Println("μ's! Muuuusic, start!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	d.Close()
}

// NewRouter
// Constructor for Router struct
func NewRouter() *Router {
	r := make(map[string]handlers.Handler)
	router := &Router{
		routes: r,
	}
	return router
}

func init() {
	// Get environment variable for discord token
	t, err := os.LookupEnv("TOKEN")
	if !err {
		log.Println("Please set env TOKEN=[AUTH_TOKEN]!")
		os.Exit(1)
	}

	// Set the token
	token = t
}

func main() {
	// Open a websocket connection to Discord
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Ara Ara:", err)
	}
	err = dg.Open()
	if err != nil {
		log.Printf("Ara Ara:", err)
		os.Exit(1)
	}

	// Generate command structs
	help := handlers.NewHelp(prefix)
	ping := handlers.NewPing()
	washi := handlers.NewWashi()
	junai := handlers.NewJunai()
	stop := handlers.NewStop()

	// Route messages based on their command
	r := NewRouter()
	r.AddHandler(help.Name, help)
	r.AddHandler(ping.Name, ping)
	r.AddHandler(washi.Name, washi)
	r.AddHandler(junai.Name, junai)
	r.AddHandler(stop.Name, stop)

	// Add descriptions to help
	help.AddDesc(&r.routes)
	r.Run(dg)
}
