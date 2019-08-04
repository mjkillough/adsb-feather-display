package ws

import (
	"log"
	"time"
)

type Fetcher interface {
	Fetch() ([]byte, error)
}

type Reporter interface {
	Report(error)
}

type Hub struct {
	fetcher  Fetcher
	reporter Reporter

	// Period to poll OpenSky API
	pollInterval time.Duration

	// Registered clients.
	clients map[*client]bool
	// Register requests from the clients.
	register chan *client
	// Unregister requests from clients.
	unregister chan *client
}

func New(fetcher Fetcher, reporter Reporter, pollInterval time.Duration) *Hub {
	return &Hub{
		fetcher:  fetcher,
		reporter: reporter,

		pollInterval: pollInterval,

		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func (h *Hub) Run() {
	timer := time.NewTimer(h.pollInterval)

	for {
		select {
		case client := <-h.register:
			log.Println("client registered")
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("client unregistered")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case <-timer.C:
			// Fetch only if we have clients listening.
			if len(h.clients) > 0 {
				bytes, err := h.fetcher.Fetch()
				if err != nil {
					log.Printf("error: %v\n", err)
				} else {
					h.broadcast(bytes)
				}
			}

			// Manually reset timer to account for slow requests to OpenSky.
			timer = time.NewTimer(h.pollInterval)
		}
	}
}

func (h *Hub) broadcast(msg []byte) {
	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			log.Println("client closed")
			close(client.send)
			delete(h.clients, client)
		}
	}
}
