// Main of Nozobot
// Modified from github.com/bwmarrin/disgord

package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VineBalloon/nozobot/client"
	"github.com/VineBalloon/nozobot/detectors"
	"github.com/VineBalloon/nozobot/handlers"
	"github.com/VineBalloon/nozobot/utils"
	"github.com/bwmarrin/discordgo"
)

// Global vars
var (
	token  string
	prefix = "><"
	router *Router
	detect *Detect
	dg     *discordgo.Session
)

// Detect Contains the detectors that will be notified on message
type Detect struct {
	detecting []detectors.Detector
}

// AddDetector Adds a detector to the detecting slice
func (d *Detect) AddDetector(detector detectors.Detector) {
	err := detector.Apiget()
	if err != nil {
		log.Fatal(err)
	}
	d.detecting = append(d.detecting, detector)
}

// Notify Notifies all detectors
func (d *Detect) Notify(cs *client.ClientState) {
	for _, det := range d.detecting {
		err := det.MsgDetect(cs)
		// Silently fail, log the error
		if err != nil {
			log.Println(err)
		}
	}
}

// NewDetect Constructor for Detect
func NewDetect() *Detect {
	return &Detect{[]detectors.Detector{}}
}

// Router Holds string->handler mappings
type Router struct {
	routes map[string]handlers.Handler
}

// AddHandler Add a string->handler mapping to the router
func (r *Router) AddHandler(handler handlers.Handler) {
	r.routes[strings.ToLower(handler.Name())] = handler
}

// Route Routes to handler from string
func (r *Router) Route(cmd string) (handlers.Handler, bool) {
	routes := r.routes[cmd]
	if routes == nil {
		return nil, false
	}
	return routes, true
}

// Run Registers our event handler(s) to the discordgo session
func (r *Router) Run(d *discordgo.Session) {
	cs := client.NewClientState()

	// Set listening
	// TODO: UpdateStatusComplex
	d.UpdateListeningStatus("you 💜")

	// Handle MessageCreate event
	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		// Update client state
		cs.UpdateState(s, m.Message)

		message := strings.TrimSpace(m.Content)

		// Send message to detectors
		detect.Notify(cs)

		if !strings.HasPrefix(message, prefix) {
			return
		}

		command := strings.ToLower(strings.TrimLeft(message, prefix))
		//log.Printf("Received message {%s}\n", message)
		args := strings.SplitN(command, " ", 2)

		handler, found := r.Route(args[0])
		if !found {
			s.ChannelMessageSend(m.ChannelID,
				utils.Italics("Ara Ara: Unknown command, see "+utils.Code(prefix+"help")))
			return
		}

		// Check roles
		member, err := s.State.Member(m.GuildID, m.Author.ID)
		if err != nil {
			member, err = s.GuildMember(m.GuildID, m.Author.ID)
			if err != nil {
				log.Println(err)
				return
			}
		}

		hroles := handler.Roles()
		groles, err := s.GuildRoles(m.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		rolesrequired := []string{}
		for hr := range hroles {
			for gr := range groles {
				if strings.ToLower(groles[gr].Name) == hroles[hr] {
					rolesrequired = append(rolesrequired, groles[gr].ID)
				}
			}
		}
		if len(hroles) > 0 {
			has := false
			mroles := member.Roles
			for _, rr := range rolesrequired {
				if utils.Stringinslice(rr, mroles) {
					has = true
				}
			}

			if !has {
				out := "Ara Ara: you need to be a " + utils.Code(hroles[0])
				others := hroles[1:]
				for hr := range others {
					out += " or a " + utils.Code(others[hr])
				}
				s.ChannelMessageSend(m.ChannelID, utils.Italics(out))
				return
			}
		}

		cs.UpdateState(s, m.Message)
		s.ChannelTyping(m.ChannelID)
		err = handler.MsgHandle(cs)
		if err != nil {
			log.Println(err)
			var out strings.Builder
			out.WriteString("Ara Ara:")
			out.WriteString(strings.Join(strings.Split(err.Error(), ":")[1:], ":"))
			s.ChannelMessageSend(m.ChannelID, utils.Italics(out.String()))
			return
		}
	})

	/* Add more event handlers here if needed */

	// Don't close the connection, wait for a kill signal
	log.Println("μ's! Muuuuusic, start!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	sig := <-sc
	log.Println("\nReceived Signal: " + sig.String())
	log.Println("Arigato, Minna-san! Sayonara!")
	d.Close()
}

// NewRouter Constructor for Router struct
func NewRouter() *Router {
	return &Router{make(map[string]handlers.Handler)}
}

// init for discordgo things
func init() {
	token, exists := os.LookupEnv("CARUDO")
	if !exists {
		log.Fatal("Missing Discord API Key: CARUDO")
	}

	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Ara Ara:", err)
	}

	err = dg.Open()
	if err != nil {
		log.Println("Ara Ara: ", err)
		os.Exit(1)
	}
}

// init for router and detector
func init() {
	// Create router, add handlers
	router = NewRouter()
	//router.AddHandler(handlers.NewGay())
	// ^DEPRECATED: Using detector instead
	router.AddHandler(handlers.NewImit())
	router.AddHandler(handlers.NewStat())
	router.AddHandler(handlers.NewPing())
	router.AddHandler(handlers.NewTarot())
	// TODO
	//router.AddHandler(handlers.NewWashi())
	//router.AddHandler(handlers.NewJunai())
	//router.AddHandler(handlers.NewStop())
	//router.AddHandler(handlers.NewLeave())

	// Add Help last to enter descriptions properly
	router.AddHandler(handlers.NewHelp(prefix, &router.routes))

	// Create detect, add detectors
	detect = NewDetect()
	detect.AddDetector(detectors.NewGay())
}

func main() {
	router.Run(dg)
}
