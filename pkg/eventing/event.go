package eventing

import (
	"fmt"
	"log"
	"time"

	"github.com/ugjka/messenger"
)

type Event struct {
	CreatedAt   time.Time
	Description string
}

type EventHandler struct {
	messenger *messenger.Messenger
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		messenger: messenger.New(0, true),
	}
}

func (h *EventHandler) Emit(s string, v ...interface{}) {
	// Construct the description
	msg := fmt.Sprintf(s, v...)

	// Log the emitted description
	log.Println(msg)

	// Broadcast the description
	h.messenger.Broadcast(Event{
		CreatedAt:   time.Now(),
		Description: msg,
	})

	h.messenger.Sub()
}

// Proxy to internal messenger
func (h *EventHandler) Sub() (client chan interface{}, err error) {
	return h.messenger.Sub()
}

// Proxy to internal messenger
func (h *EventHandler) Unsub(client chan interface{}) {
	h.messenger.Unsub(client)
}
