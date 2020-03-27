package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func (serv *Server) ApplyHandlers() {
	serv.router.Handle("/*", http.FileServer(http.Dir("./web")))
	serv.router.Get("/socket", serv.socketHandler)
}

func (serv *Server) socketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := serv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("websocket err: %v", err)
	}

	go func() {
		for {
			<-time.After(5 * time.Second)
			msg := Message{
				Type: MTPing,
			}
			if err := ws.WriteJSON(msg); err != nil {
				log.Printf("ws send ping err: %v", err)
			}
		}
	}()

	id := uuid.New().String()
	serv.submutex.Lock()
	serv.subscribers[id] = func(msg string) error {
		m := Message{
			Type: MTMessage,
			Data: msg,
		}
		if err := ws.WriteJSON(m); err != nil {
			log.Printf("ws msg fetch err: %v", err)
		}
		return nil
	}
	serv.submutex.Unlock()

	for {
		msg := Message{}
		if err := ws.ReadJSON(&msg); err != nil {
			if !websocket.IsCloseError(err, 1001) {
				log.Fatalf("ws msg read err: %v", err)
			}
			break
		}

		if msg.Type == MTPong {
			continue
		}

		if msg.Type == MTMessage {
			fmt.Println(msg.Data)
			serv.submutex.Lock()
			for _, sub := range serv.subscribers {
				if err := sub(msg.Data); err != nil {
					log.Fatalf("ws msg subs err: %v", err)
				}
			}
			serv.submutex.Unlock()
		}
	}

	fmt.Println("CLSOED")
	defer func() {
		serv.submutex.Lock()
		delete(serv.subscribers, id)
		serv.submutex.Unlock()
	}()
}
