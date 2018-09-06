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
	"time"

	"github.com/bwmarrin/discordgo"
)

// Global var for our token
var (
	Token  string
	Prefix = "!"
)

// MHandler type for message-invoked commands
type MHandler func(s *discordgo.Session, m *discordgo.MessageCreate) error

// Handler struct
type Handler struct {
	mhandler MHandler
	desc     string
	roles    []string
	channels []string
}

// Setter for description
func (h *Handler) Desc(d string) {
	h.desc = d
}

// Setter for roles
func (h *Handler) Roles(r []string) {
	h.roles = r
}

// Setter for channels
func (h *Handler) Channels(c []string) {
	h.channels = c
}

// Simple constructor for Handler
func NewHandler(m MHandler) *Handler {
	return &Handler{
		mhandler: m,
		desc:     "",
		roles:    nil,
		channels: nil,
	}
}

// Helper function to get guild and vc id from state and message
// Pass into constructor as `NewVoiceRoom(VoiceInfoFromMessage(session, user))
func VoiceInfoFromMessage(s *discordgo.Session, m *discordgo.Message) (string, string) {
	// Get the guild and guild ID
	mchannel, _ := s.Channel(m.ChannelID)
	guildId := mchannel.GuildID
	guild, _ := s.Guild(guildId)

	// Get channel id
	u := m.Author
	channel := ""
	for _, vs := range guild.VoiceStates {
		if vs.UserID == u.ID {
			channel = vs.ChannelID
		}
	}

	return guildId, channel
}

type VoiceRoom struct {
	guild      string
	id         string
	Connection *discordgo.VoiceConnection
}

// Construct a new voice room
func NewVoiceRoom(guild, channel string) (*VoiceRoom, error) {
	if channel == "" {
		return nil, errors.New("nozobot: user not in voice channel")
	}

	return &VoiceRoom{
		guild:      guild,
		id:         channel,
		Connection: nil,
	}, nil
}

// Connect to the voice channel
func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.guild, v.id, false, false)
	v.Connection = vc
	return err
}

// Close the voice connection
func (v *VoiceRoom) Close() {
	v.Connection.Disconnect()
	v.Connection = nil
}

// Router struct to hold our string->router mappings
type Router struct {
	routes map[string]Handler
}

// Router method to add a string->handler mapping
func (r *Router) AddHandler(cmd string, handler Handler) {
	r.routes[cmd] = handler
}

// Method to run the router
func (r *Router) Run(d *discordgo.Session) {
	// Add anonymous function to route messages to handlers
	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Don't talk to yourself
		if m.Author.ID == s.State.User.ID {
			return
		}

		message := strings.TrimSpace(m.Content)
		if !strings.HasPrefix(message, Prefix) {
			return
		}

		// Parse message
		command := strings.ToLower(strings.TrimLeft(message, Prefix))

		//fmt.Printf("Received message {%s}\n", message)

		// Split the command into 2 substrings
		args := strings.SplitN(command, " ", 2)

		// Route message to the handler
		handler, found := r.routes[args[0]]
		if !found {
			err := "Unknown command, type " + Prefix + "help"
			s.ChannelMessageSend(m.ChannelID, err)
			return
		}

		// TODO make !help default handler that iterates through handlers
		// and make handlers a struct actually that have a description

		// Call mhandler
		err := handler.mhandler(s, m)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
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
	r := make(map[string]Handler)
	// Help is a default handler
	router := &Router{
		routes: r,
	}
	router.AddHandler("help", *NewHandler(HandleHelp))
	return router
}

// Default help command
func HandleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// TODO Iterate through the router's handlers
	_, err := s.ChannelMessageSend(m.ChannelID, "Still under construction!")
	if err != nil {
		return err
	}
	return nil
}

/************************/
/********COMMANDS********/

// TODO Make each handler a struct that implements the handler interface
// ping
func HandlePing(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}

// TODO make nozomi join VC and play the audio
// washi
func HandleWashi(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Washi Washi!")
	if err != nil {
		return err
	}

	// Create a voice room
	vr, err := NewVoiceRoom(VoiceInfoFromMessage(s, m.Message))
	if err != nil {
		return err
	}

	// Connect to the voice room
	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Sleep for 5 seconds
	time.Sleep(time.Second * 5)

	// Close the voice connection
	vr.Close()
	return nil
}

// tarot
func HandleTarot(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	// Route messages based on their command
	r := NewRouter()
	r.AddHandler("ping", *NewHandler(HandlePing))
	r.AddHandler("washi", *NewHandler(HandleWashi))
	r.Run(dg)
}
