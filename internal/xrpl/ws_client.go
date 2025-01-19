package xrpl

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Connection *websocket.Conn
}

// create a new WebSocket client
func NewWebSocketClient(url string) (*WebSocketClient, error) {
	// create a new WebSocket dialer with a 10 second timeout to establish a connection 
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(url, nil) // establish a connection to the WebSocket server
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to WebSocket server: %s", url)
	return &WebSocketClient{Connection: conn}, nil // return a new WebSocket client
}

// send a request to the WebSocket server
func (wsc *WebSocketClient) Subscribe(request interface{}) error {
	reqJSON, err := json.Marshal(request) // convert the request to JSON
	if err != nil {
		return err
	}
	// send the request to the WebSocket server as a text message 
	return wsc.Connection.WriteMessage(websocket.TextMessage, reqJSON)
}

// read messages from the WebSocket server
func (wsc *WebSocketClient) ReadMessages(callback func(message []byte)) {
	defer wsc.Connection.Close() // close the connection when the function returns

	for {
    _, msg, err := wsc.Connection.ReadMessage()
    if err != nil {
        log.Printf("⚠️ WebSocket desconectado. Tentando reconectar... Erro: %v", err)
        time.Sleep(5 * time.Second) // Esperar antes de tentar reconectar
        wsc.Connection, _, err = websocket.DefaultDialer.Dial(wsc.Connection.RemoteAddr().String(), nil)
        if err != nil {
            log.Printf("❌ Falha ao reconectar: %v", err)
            continue
        }
        log.Println("✅ Reconexão bem-sucedida!")
        continue
    }
    callback(msg)
	}
}
