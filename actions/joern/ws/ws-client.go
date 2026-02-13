package wsHandler

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/SorenHQ/joern-port/etc"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)


var rh *ResultHandlers
type ResultHandlers struct {
	// Define your result handlers here
	conn           *websocket.Conn
	serverUrl      string
	messageChan chan string
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
	response, statusCode, _ := etc.CustomCall(ctx, "GET", nil, url, nil)
	if statusCode!=200{
		fmt.Printf("joern server response code is %d\n",statusCode)
		rh.messageChan <- string(message)+"||"+string("service not available")
		return
	}

	jsonObj, err := etc.ParseJoernStdoutJsonObject(response)
	if err != nil {
		log.Println("Error parsing response:", err)
		rh.messageChan <- string(message)+"||"+string(response)
		return
	}
	// Marshal the JSON object back to string for the channel
	jsonBytes, err := sonic.Marshal(jsonObj)
	if err != nil {
		log.Println("Error marshaling JSON object:", err)
		rh.messageChan <- string(message)+"||"+string(response)
		return
	}
	rh.messageChan <- string(message)+"||"+string(jsonBytes)
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
func (rh *ResultHandlers) GetResultChannel() chan string {
	return  rh.messageChan
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

func NewJoernResultHandlers(serverURL string) (*ResultHandlers, error) {
	url, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	rh = &ResultHandlers{
		conn:           nil,
		serverUrl:      url.Host,
		messageChan: make(chan string),

	}
	err = rh.connectToServer()
	if err != nil {
		return nil, err
	}
	return rh, nil
}

func GetJoernResultHandler()*ResultHandlers{
	if rh ==nil{
		return nil
	}

	return  rh
}