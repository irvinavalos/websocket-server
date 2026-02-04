package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	WSPort = ":8080"
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

func handleWS(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/", handleWS)

	log.Fatal(http.ListenAndServe(WSPort, nil))
}
