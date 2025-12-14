package server

import (
	"fmt"
	"net/http"

	"plan/internal/events"
)

func handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if flushing is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Create client channel
	client := make(chan events.Event, 10)
	events.DefaultBroker.Register(client)
	defer events.DefaultBroker.Unregister(client)

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	// Stream events
	for {
		select {
		case event := <-client:
			sse, err := event.FormatSSE()
			if err != nil {
				continue
			}
			fmt.Fprint(w, sse)
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}
