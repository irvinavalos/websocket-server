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
	WSPort     = ":8000"
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
	Clients       map[string]*Client
	JoinServerCh  chan *Client
	LeaveServerCh chan *Client
	mu            *sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Clients:       map[string]*Client{},
		JoinServerCh:  make(chan *Client, 64),
		LeaveServerCh: make(chan *Client, 64),
		mu:            new(sync.RWMutex),
	}
}

func (s *Server) clientJoiningServer(c *Client) {
	s.Clients[c.ID] = c
	log.Printf("Client joining...\nClientID: %s\n", c.ID)
}

func (s *Server) clientLeavingServer(c *Client) {
	delete(s.Clients, c.ID)
	log.Printf("Client leaving...\nClientID: %s\n", c.ID)
}

func (s *Server) AcceptLoop() {
	for {
		select {
		case c := <-s.JoinServerCh:
			s.clientJoiningServer(c)
		case c := <-s.LeaveServerCh:
			s.clientLeavingServer(c)
		}
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
	s.JoinServerCh <- client
}

func wsServerStart() {
	s := NewServer()
	go s.AcceptLoop()
	http.HandleFunc("/", s.handleWS)
	log.Println("Starting server...")

	log.Fatal(http.ListenAndServe(WSPort, nil))
}

func main() {
	wsServerStart()
}
