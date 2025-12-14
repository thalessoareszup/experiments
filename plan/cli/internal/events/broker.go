package events

import (
	"encoding/json"
	"sync"
)

type EventType string

const (
	PlanCreated EventType = "plan:created"
	PlanUpdated EventType = "plan:updated"
	PlanDeleted EventType = "plan:deleted"
	StepCreated EventType = "step:created"
	StepUpdated EventType = "step:updated"
	StepDeleted EventType = "step:deleted"
)

type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

type Broker struct {
	clients    map[chan Event]bool
	register   chan chan Event
	unregister chan chan Event
	broadcast  chan Event
	mu         sync.RWMutex
}

var DefaultBroker *Broker

func init() {
	DefaultBroker = NewBroker()
	go DefaultBroker.Run()
}

func NewBroker() *Broker {
	return &Broker{
		clients:    make(map[chan Event]bool),
		register:   make(chan chan Event),
		unregister: make(chan chan Event),
		broadcast:  make(chan Event, 100),
	}
}

func (b *Broker) Run() {
	for {
		select {
		case client := <-b.register:
			b.mu.Lock()
			b.clients[client] = true
			b.mu.Unlock()

		case client := <-b.unregister:
			b.mu.Lock()
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client)
			}
			b.mu.Unlock()

		case event := <-b.broadcast:
			b.mu.RLock()
			for client := range b.clients {
				select {
				case client <- event:
				default:
					// Client buffer full, skip
				}
			}
			b.mu.RUnlock()
		}
	}
}

func (b *Broker) Register(client chan Event) {
	b.register <- client
}

func (b *Broker) Unregister(client chan Event) {
	b.unregister <- client
}

func (b *Broker) Broadcast(event Event) {
	b.broadcast <- event
}

// Emit is a convenience function to broadcast an event
func Emit(eventType EventType, data interface{}) {
	DefaultBroker.Broadcast(Event{Type: eventType, Data: data})
}

// FormatSSE formats an event for SSE transmission
func (e Event) FormatSSE() (string, error) {
	data, err := json.Marshal(e.Data)
	if err != nil {
		return "", err
	}
	return "event: " + string(e.Type) + "\ndata: " + string(data) + "\n\n", nil
}
