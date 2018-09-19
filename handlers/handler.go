package handlers

import (
	"github.com/VineBalloon/nozobot/client"
)

// Handler interface for commands to implement
type Handler interface {
	Desc() string
	Roles() []string
	Channels() []string
	Handle(*client.ClientState) error
}
