package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	BufferSize = 512
	WSPort     = ":8080"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	mu   *sync.RWMutex
}

func NewClient(conn *websocket.Conn) *Client {
	id := rand.Text()[:10]
	return &Client{
		ID:   id,
		Conn: conn,
		mu:   new(sync.RWMutex),
	}
}

type Server struct {
	Clients []*Client
	mu      *sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Clients: []*Client{},
		mu:      new(sync.RWMutex),
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to WS
	upgrader := websocket.Upgrader{
		ReadBufferSize:  BufferSize,
		WriteBufferSize: BufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading HTTP connection to WS: %s\n", err)
		return
	}

	// Create client and add it to server
	client := NewClient(conn)
	s.Clients = append(s.Clients, client) // FIX: race condition
}

func main() {
	s := NewServer()
	http.HandleFunc("/", s.handleWS)

	log.Fatal(http.ListenAndServe(WSPort, nil))
}
