package wsHandler

import (
	"context"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func wsConnection(serverURL string, messageChan chan []byte) error {
	// Replace with your WebSocket server URL
	if conn != nil {
		conn.Close()
		log.Println("reconnect to ", serverURL)
	}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
		return err
	}

	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			messageChan <- message
		}

	}()
	return nil
}

func ConnectToServer(serverURL string, messageHandler func(serverURL string, message []byte)) error {
	messageChan := make(chan []byte)
	ctx := context.WithoutCancel(context.Background())
	wsUrl := strings.ReplaceAll(serverURL, "http://", "ws://")
	err := wsConnection(wsUrl+"/connect", messageChan)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
		return err
	}
	go func() {
		for {
			select {
			case message := <-messageChan:
				messageHandler(serverURL, message)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
