package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Appointat/Responsive-AI-Clusters-in-Supply-Chain/agent"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	Warehouse *agent.Warehouse
	Outlets   []*agent.Outlet
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	mutex     sync.Mutex
}

func NewServer(wh *agent.Warehouse, outlets []*agent.Outlet) *Server {
	return &Server{
		Warehouse: wh,
		Outlets:   outlets,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	defer ws.Close()

	s.mutex.Lock()
	s.Clients[ws] = true
	s.mutex.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			s.mutex.Lock()
			delete(s.Clients, ws)
			s.mutex.Unlock()
			break
		}
	}
}

func (s *Server) HandleMessages() {
	for {
		msg := <-s.Broadcast
		s.mutex.Lock()
		for client := range s.Clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(s.Clients, client)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *Server) SendUpdate(data interface{}) {
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		log.Printf("Json marshal error: %v", err)
		return
	}
	s.Broadcast <- jsonMsg
}
