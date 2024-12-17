package account_channels

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Connection *websocket.Conn
}

func NewWebSocketClient(url string) (*WebSocketClient, error) {

	dialer:= websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}
	log.Println("Connected to WebSocket Server")
	return &WebSocketClient{Connection: conn}, nil
}

func (wsc *WebSocketClient) SendResquest(request inferface{}) (*XRPAccountChannelsResponse, error) {
	
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	
	err = wsc.Connection.WriteMessage(websocket.TextMessage, reqJSON)
	if err != nil {
		return nil, err
	}

	_, message, err := wsc.Connection.ReadMessage()
	if err != nil {
		return nil, err
	}

	var result XRPAccountChannelsResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
} 