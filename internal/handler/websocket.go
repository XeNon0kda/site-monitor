package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка апгрейда WebSocket: %v", err)
		return
	}
	defer conn.Close()

	statusCh := h.monitor.Subscribe()
	defer h.monitor.Unsubscribe(statusCh)

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()

	for {
		select {
		case status, ok := <-statusCh:
			if !ok {
				return
			}
			if err := conn.WriteJSON(status); err != nil {
				log.Printf("Ошибка записи в WebSocket: %v", err)
				return
			}
		}
	}
}