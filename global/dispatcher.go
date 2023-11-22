package global

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	Clients    map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Locker     sync.RWMutex
}

type Message struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Msg  string `json:"msg"`
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[int]*Client),
		Locker:     sync.RWMutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Locker.Lock()
			h.Clients[client.UserId] = client
			h.Locker.Unlock()
		case client := <-h.Unregister:
			h.Locker.Lock()
			if _, ok := h.Clients[client.UserId]; ok {
				delete(h.Clients, client.UserId)
				close(client.Send)
			}
			h.Locker.Unlock()
		}
	}
}

func (h *Hub) Dispatch(userId int, message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil || msg.From != userId {
		return
	}
	h.Locker.RLock()
	if client, ok := h.Clients[msg.To]; ok {
		select {
		case client.Send <- message:
			// Save msg into database, setting has_sent true
			_, err := Database.Exec(`INSERT INTO "message" (from, to, message, has_sent, time) VALUES ($1, $2, $3, true, $4)`, msg.From, msg.To, msg.Msg, time.Now())
			if err != nil {
				log.Fatalf("Save message into database failed: %v", err)
				return
			}
		default:
			// Save msg into database, setting has_sent false
			_, err := Database.Exec(`INSERT INTO "message" (from, to, message, has_sent, time) VALUES ($1, $2, $3, false, $4)`, msg.From, msg.To, msg.Msg, time.Now())
			if err != nil {
				log.Fatalf("Save message into database failed: %v", err)
				return
			}
			h.Unregister <- client
		}
	}
	h.Locker.RUnlock()
}