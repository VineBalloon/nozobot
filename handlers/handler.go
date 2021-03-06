package handlers

import "github.com/VineBalloon/nozobot/client"

// Handler
// The interface for all commands
type Handler interface {
	Name() string       /* Name returns the name of the command */
	Desc() string       /* Desc returns the description of the command */
	Roles() []string    /* Roles returns the permitted roles for the command */
	Channels() []string /* Channels returns the permitted channels */

	MsgHandle(*client.ClientState) error /* Handle handles message events */
	// Add more event handlers here
}

// Initialise
func init() {

}
