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

// Handler interface for commands to implement
type Handler interface {
	Desc() string
	Roles() []string
	Channels() []string
	Handle(*discordgo.Session, *discordgo.MessageCreate) error
}

// MHandler type for message-invoked commands
type MHandler func(s *discordgo.Session, m *discordgo.MessageCreate) error

// Voice Room Constructor
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

type VoiceRoom struct {
	guild      string
	id         string
	Connection *discordgo.VoiceConnection
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

func (v *VoiceRoom) Connect(s *discordgo.Session) error {
	// Attempt to generate a voice connection
	vc, err := s.ChannelVoiceJoin(v.guild, v.id, false, false)
	v.Connection = vc
	return err
}

func (v *VoiceRoom) Close() {
	v.Connection.Disconnect()
	v.Connection = nil
}

// Router struct to hold our string->router mappings
type Router struct {
	routes map[string]Handler
}

func (r *Router) AddHandler(cmd string, handler Handler) {
	r.routes[strings.ToLower(cmd)] = handler
}

func (r *Router) Route(cmd string) (Handler, bool) {
	routes := r.routes[cmd]
	if routes == nil {
		return nil, false
	}
	return routes, true
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
		handler, found := r.Route(args[0])
		if !found {
			err := "Unknown command, use " + Prefix + "help"
			s.ChannelMessageSend(m.ChannelID, err)
			return
		}

		//handler := *handlerp
		// Call handler method
		err := handler.Handle(s, m)

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
	router := &Router{
		routes: r,
	}
	return router
}

// Default help command
func NewHelp(n string) *Help {
	return &Help{
		"help",
		nil,
	}
}

type Help struct {
	name         string
	descriptions map[string]string
}

func (h *Help) AddDesc(r *map[string]Handler) {
	h.descriptions = make(map[string]string)
	for cmd, handler := range *r {
		h.descriptions[cmd] = handler.Desc()
	}
	fmt.Println(h.descriptions)
}

// TODO Iterate through the router's handlers and refactor
func (h *Help) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
	out := "Commands:\n"
	for name, desc := range h.descriptions {
		out += name + ": " + desc + "\n"
	}
	_, err := s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		return err
	}
	return nil
}

func (h *Help) Desc() string {
	return "Nozomi helps you write out this command!"
}

func (h *Help) Roles() []string {
	return nil
}

func (h *Help) Channels() []string {
	return nil
}

/************************/
/********COMMANDS********/

func NewPing(n string) *Ping {
	return &Ping{
		"Ping",
	}
}

func NewWashi(n string) *Washi {
	return &Washi{
		"Washi",
	}
}

func NewTarot(n string) *Tarot {
	return &Tarot{
		"Tarot",
	}
}

/*----------------------*/
// TODO put each command into their own file

type Ping struct {
	name string
}

func (p *Ping) Desc() string {
	return "Ping pong with Nozomi :ping_pong:"
}

func (p *Ping) Roles() []string {
	return nil
}

func (p *Ping) Channels() []string {
	return nil
}

func (p *Ping) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}

type Washi struct {
	name string
}

func (w *Washi) Desc() string {
	return "Nozomi's washi washi will follow you into Voice as well!"
}

func (w *Washi) Roles() []string {
	return nil
}

func (w *Washi) Channels() []string {
	return nil
}

func (w *Washi) Handle(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Washi Washi!")
	if err != nil {
		return err
	}

	// Attempt to join a voice room
	vr, err := NewVoiceRoom(VoiceInfoFromMessage(s, m.Message))
	if err != nil {
		return err
	}

	err = vr.Connect(s)
	if err != nil {
		return err
	}

	// Sleep for 5 seconds
	// TODO make nozomi play the audio
	time.Sleep(time.Second * 5)

	// Close the voice connection
	vr.Close()
	return nil
}

type Tarot struct {
	name string
}

func (t *Tarot) Desc() string {
	return "Nozomi decides your fate!"
}

func (t *Tarot) Roles() []string {
	return nil
}

func (t *Tarot) Channels() []string {
	return nil
}

func (t *Tarot) HandleTarot(s *discordgo.Session, m *discordgo.MessageCreate) error {
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

	// Genearte command structs
	help := NewHelp("help")
	ping := NewPing("ping")
	washi := NewWashi("washi")

	// Route messages based on their command
	r := NewRouter()
	r.AddHandler(help.name, help)
	r.AddHandler(ping.name, ping)
	r.AddHandler(washi.name, washi)

	// Add help messages to help
	help.AddDesc(&r.routes)
	fmt.Println(help.descriptions)
	r.Run(dg)
}
