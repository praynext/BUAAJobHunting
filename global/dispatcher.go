package global

import (
	"BUAAJobHunting/model"
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
	Time string `json:"time"`
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
	msg.Time = time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006/01/02 15:04:05")
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Marshal messages failed: %v", err)
		return
	}
	h.Locker.RLock()
	if client, ok := h.Clients[msg.To]; ok {
		select {
		case client.Send <- byteMsg:
			// Save msg into database, setting has_sent true
			_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, true, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
			if err != nil {
				log.Fatalf("Save message into database failed: %v", err)
				return
			}
		default:
			// Save msg into database, setting has_sent false
			_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, false, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
			if err != nil {
				log.Fatalf("Save message into database failed: %v", err)
				return
			}
			h.Unregister <- client
		}
	} else {
		// Save msg into database, setting has_sent false
		_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, false, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
		if err != nil {
			log.Fatalf("Save message into database failed: %v", err)
			return
		}
	}
	h.Locker.RUnlock()
	if _, err = Database.Exec(`DELETE FROM last_chat WHERE "from" = $1 AND "to" = $2`, msg.To, msg.From); err != nil {
		log.Fatalf("Update last_chat failed: %v", err)
		return
	}
	if _, err = Database.Exec(`INSERT INTO last_chat ("from", "to", time) VALUES ($1, $2, $3) 
        ON CONFLICT ("from", "to") DO UPDATE set time = $4`, msg.From, msg.To, msg.Time, msg.Time); err != nil {
		log.Fatalf("Update last_chat failed: %v", err)
		return
	}
}

func (h *Hub) Remind(lastChats []model.LastChat) {
	h.Locker.RLock()
	for _, lastChat := range lastChats {
		msg := Message{
			From: 0,
			To:   lastChat.To,
			Msg:  "您有未读消息，请及时查看",
			Time: time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006/01/02 15:04:05"),
		}
		byteMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("Marshal messages failed: %v", err)
			return
		}
		if client, ok := h.Clients[msg.To]; ok {
			select {
			case client.Send <- byteMsg:
				// Save msg into database, setting has_sent true
				_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, true, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
				if err != nil {
					log.Fatalf("Save message into database failed: %v", err)
					return
				}
			default:
				// Save msg into database, setting has_sent false
				_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, false, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
				if err != nil {
					log.Fatalf("Save message into database failed: %v", err)
					return
				}
				h.Unregister <- client
			}
		} else {
			// Save msg into database, setting has_sent false
			_, err := Database.Exec(`INSERT INTO "message" ("from", "to", message, has_sent, time) VALUES ($1, $2, $3, false, $4)`, msg.From, msg.To, msg.Msg, msg.Time)
			if err != nil {
				log.Fatalf("Save message into database failed: %v", err)
				return
			}
		}
	}
	h.Locker.RUnlock()
	for _, lastChat := range lastChats {
		if _, err := Database.Exec(`DELETE FROM last_chat WHERE "from" = $1 AND "to" = $2`, lastChat.From, lastChat.To); err != nil {
			log.Fatalf("Update last_chat failed: %v", err)
			return
		}
	}
}
