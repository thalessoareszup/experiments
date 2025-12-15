package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"plan/internal/events"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (CORS is handled at middleware level)
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WSClient represents a connected WebSocket client
type WSClient struct {
	conn      *websocket.Conn
	send      chan []byte
	eventChan chan events.Event
}

// handleWebSocket upgrades HTTP connection to WebSocket and manages the connection lifecycle
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	client := &WSClient{
		conn:      conn,
		send:      make(chan []byte, 256),
		eventChan: make(chan events.Event, 10),
	}

	// Register client with event broker
	events.DefaultBroker.Register(client.eventChan)
	defer func() {
		events.DefaultBroker.Unregister(client.eventChan)
		client.conn.Close()
	}()

	// Send initial connected message
	if err := client.sendMessage(&events.WSMessage{
		Type:      "connected",
		Data:      json.RawMessage(`{"status":"connected"}`),
		Timestamp: time.Now(),
		MessageID: events.GenerateMessageID(),
	}); err != nil {
		log.Printf("Failed to send initial message: %v", err)
		return
	}

	// Start read and write goroutines
	go client.readLoop()
	go client.writeLoop()

	// Main event loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-client.eventChan:
			if !ok {
				return
			}
			// Format event as WebSocket message
			msg, err := event.FormatWebSocket()
			if err != nil {
				log.Printf("Failed to format event: %v", err)
				continue
			}

			// Send event
			if err := client.sendMessage(msg); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			if err := client.conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(10*time.Second),
			); err != nil {
				log.Printf("Ping error: %v", err)
				return
			}
		}
	}
}

// readLoop reads messages from the client (for future bidirectional features)
func (c *WSClient) readLoop() {
	defer c.conn.Close()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg map[string]interface{}
		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		// Handle client messages (future feature for bidirectional communication)
		// For now, just log them
		log.Printf("Received message from client: %v", msg)
	}
}

// writeLoop writes messages from the send channel to the WebSocket
func (c *WSClient) writeLoop() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// sendMessage sends a WebSocket message to the client
func (c *WSClient) sendMessage(msg *events.WSMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case c.send <- data:
		return nil
	default:
		// Buffer full, close connection to avoid memory buildup
		log.Printf("Client send buffer full, closing connection")
		return fmt.Errorf("send buffer full")
	}
}

