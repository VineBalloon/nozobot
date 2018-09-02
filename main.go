// Main of Nozobot
// Modified from github.com/bwmarrin/disgord

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Global var for our token
var Token string

// Handler type for commands
type Handler func(msg string) error

// Router struct to hold our string->router mappings
type Router struct {
	routes map[string]Handler
}

// Constructor for Router struct
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

// Router method to add a string->handler mapping
func (r *Router) AddHandler(cmd string, handler Handler) {
	r.routes[cmd] = handler
}

// Method to run the router
func (r *Router) Run() {
	// Open the session

}

/************************/
/********COMMANDS********/
// Ping command
func HandlePing(msg string) error {
	return nil
}

// Washi
func HandleWashi(msg string) error {
	fmt.Println("Washi Washi!")
	return nil
}

/********COMMANDS********/
/************************/

func init() {
	// Get environment variable for discord token
	token, err := os.LookupEnv("TOKEN")
	if !err {
		log.Println("Please set env TOKEN=[AUTH_TOKEN]!")
		os.Exit(1)
	}

	// Set the token
	Token = token
}

func main() {
	// Open a websocket connection to Discord
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Ara Ara:", err)
	}
	err = dg.Open()
	if err != nil {
		log.Printf("Ara Ara:", err)
		os.Exit(1)
	}

	// Add handler
	dg.AddHandler(messageCreate)

	/*
		r := NewRouter()
		r.AddHandler("ping", HandlePing)
		r.AddHandler("washi", HandleWashi)
		r.Run()
	*/

	// Wait for a kill signal
	fmt.Println("Ctrl-C to kill")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't talk to yourself
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := strings.TrimSpace(m.Content)
	//fmt.Printf("Received message {%s}\n", message)

	// Simple ping message
	if message == "ping" {
		sent, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			fmt.Println("Ara Ara:", err)
		}
		//fmt.Println("Send message:", sent)
	}
}
