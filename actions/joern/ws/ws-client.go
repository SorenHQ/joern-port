package wsHandler

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/SorenHQ/joern-port/etc"

	"github.com/gorilla/websocket"
)

type MessageHandler interface {
	Recv(string,string)
}
type ResultHandlers struct {
	// Define your result handlers here
	conn           *websocket.Conn
	serverUrl      string
	messageHandler MessageHandler
}

func (rh *ResultHandlers) getResult(message []byte) {
	// fmt.Println(string(message))
	if string(message) == "connected" {
		log.Default().Printf("websocket connected to %s\n", rh.serverUrl)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := fmt.Sprintf("http://%s/result/%s", rh.serverUrl, string(message))
	response, _, _ := etc.CustomCall(ctx, "GET", nil, url, nil)
	respBody, err := etc.ParseJoernStdoutToString(response)
	if err != nil {
		log.Println("Error parsing response:", err)
		return
	}
	rh.messageHandler.Recv(string(message),respBody)
}
func (rh *ResultHandlers) wsConnection(serverURL string, messageChan chan []byte) error {
	// Replace with your WebSocket server URL
	if rh.conn != nil {
		rh.conn.Close()
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

func (rh *ResultHandlers) connectToServer() error {
	messageChan := make(chan []byte)
	ctx := context.WithoutCancel(context.Background())
	wsUrl := fmt.Sprintf("ws://%s", rh.serverUrl)
	err := rh.wsConnection(wsUrl+"/connect", messageChan)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
		return err
	}
	go func() {
		for {
			select {
			case message := <-messageChan:
				rh.getResult(message)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func NewJoernResultHandlers(serverURL string, messHandler MessageHandler) (*ResultHandlers, error) {
	url, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	rh := ResultHandlers{
		conn:           nil,
		serverUrl:      url.Host,
		messageHandler: messHandler,
	}
	err = rh.connectToServer()
	if err != nil {
		return nil, err
	}
	return &rh, nil
}
