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
	"github.com/VineBalloon/nozobot/handlers"
	h "github.com/VineBalloon/nozobot/helpers"
	"github.com/bwmarrin/discordgo"
)

// Global vars
var (
	token  string
	prefix = "><"
	router *Router
	dg     *discordgo.Session
)

// Router
// Struct to hold our string->router mappings
type Router struct {
	routes map[string]handlers.Handler
}

// AddHandler
// Adds a string->handler mapping to the router
func (r *Router) AddHandler(handler handlers.Handler) {
	r.routes[strings.ToLower(handler.Name())] = handler
}

// Route
// Routes to handler from string
func (r *Router) Route(cmd string) (handlers.Handler, bool) {
	routes := r.routes[cmd]
	if routes == nil {
		return nil, false
	}
	return routes, true
}

// Run
// Registers our event handler(s) to the discordgo session
func (r *Router) Run(d *discordgo.Session) {
	cs := client.NewClientState()

	// Add handler to MessageCreate event
	// this gets called when a message is read by the bot
	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		message := strings.TrimSpace(m.Content)
		if !strings.HasPrefix(message, prefix) {
			return
		}

		command := strings.ToLower(strings.TrimLeft(message, prefix))
		//log.Printf("Received message {%s}\n", message)
		args := strings.SplitN(command, " ", 2)

		handler, found := r.Route(args[0])
		if !found {
			s.ChannelMessageSend(m.ChannelID,
				h.Italics("Ara Ara: Unknown command, see "+h.Code(prefix+"help")))
			return
		}

		// Check roles
		var member *discordgo.Member
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
			for rr := range rolesrequired {
				for mr := range mroles {
					if rolesrequired[rr] == mroles[mr] {
						has = true
					}
				}
			}

			if !has {
				out := "Ara Ara: you need to be a " + h.Code(hroles[0])
				others := hroles[1:]
				for hr := range others {
					out += " or a " + h.Code(others[hr])
				}
				s.ChannelMessageSend(m.ChannelID, h.Italics(out))
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
			s.ChannelMessageSend(m.ChannelID, h.Italics(out.String()))
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

// NewRouter
// Constructor for Router struct
func NewRouter() *Router {
	r := make(map[string]handlers.Handler)
	router := &Router{
		routes: r,
	}
	return router
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
		log.Printf("Ara Ara: ", err)
		os.Exit(1)
	}
}

// init was so good, I just had to make another one!
func init() {
	// Create router, register handlers
	router = NewRouter()
	router.AddHandler(handlers.NewPing())
	router.AddHandler(handlers.NewGay())
	router.AddHandler(handlers.NewImit())
	// TODO
	//router.AddHandler(handlers.NewWashi())
	//router.AddHandler(handlers.NewJunai())
	//router.AddHandler(handlers.NewStop())
	//router.AddHandler(handlers.NewLeave())

	// Add Help last to enter descriptions properly
	router.AddHandler(handlers.NewHelp(prefix, &router.routes))
}

func main() {
	router.Run(dg)
}
