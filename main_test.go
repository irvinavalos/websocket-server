package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Host = "ws://localhost"
)

type TestConfig struct {
	ClientCount int
	wg          *sync.WaitGroup
}

func DialServer(wg *sync.WaitGroup) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(fmt.Sprintf("%s%s", Host, WSPort), nil)
	if err != nil {
		log.Fatal("Failed to connect to WS ->", err)
	}

	defer func() {
		conn.Close()
		wg.Done()
	}()

	log.Println("Connecting to server... -> ", conn.LocalAddr().String())

	time.Sleep(2 * time.Second)
}

func TestConnection(t *testing.T) {
	go wsServerStart()

	time.Sleep(1 * time.Second)

	tc := TestConfig{
		ClientCount: 50,
		wg:          new(sync.WaitGroup),
	}

	tc.wg.Add(tc.ClientCount)

	for range tc.ClientCount {
		go DialServer(tc.wg)
	}

	tc.wg.Wait()

	log.Println("Exiting test")
}
