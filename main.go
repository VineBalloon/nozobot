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
	"github.com/VineBalloon/nozobot/helpers"
	"github.com/bwmarrin/discordgo"
)

// Global var for our token
var (
	token  string
	prefix = "!"
	help   = []string{}
)

// Router struct to hold our string->router mappings
type Router struct {
	routes map[string]handlers.Handler
}

func (r *Router) AddHandler(cmd string, handler handlers.Handler) {
	r.routes[strings.ToLower(cmd)] = handler
}

func (r *Router) Route(cmd string) (handlers.Handler, bool) {
	routes := r.routes[cmd]
	if routes == nil {
		return nil, false
	}
	return routes, true
}

// Method to run the router
func (r *Router) Run(d *discordgo.Session) {
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
			err := errors.New("Unknown command, use " + helpers.Code(prefix+"help"))
			s.ChannelMessageSend(m.ChannelID,
				helpers.Italics("Error: "+err.Error()))
			return
		}

		// Call handler method
		// TODO pass args
		err := handler.Handle(client.NewClientState(s, m.Message))

		if err != nil {
			s.ChannelMessageSend(m.ChannelID,
				helpers.Italics("Error: "+err.Error()))
			return
		}
	})

	// Don't close the connection, wait for a kill signal
	fmt.Println("Ctrl-C to kill")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	d.Close()
}

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

	// Genearte command structs
	help := handlers.NewHelp("help", prefix)
	ping := handlers.NewPing("ping")
	washi := handlers.NewWashi("washi")
	junai := handlers.NewJunai("junai")

	// Route messages based on their command
	r := NewRouter()
	r.AddHandler(help.Name, help)
	r.AddHandler(ping.Name, ping)
	r.AddHandler(washi.Name, washi)
	r.AddHandler(junai.Name, junai)

	// Add descriptions to help
	help.AddDesc(&r.routes)
	r.Run(dg)
}
