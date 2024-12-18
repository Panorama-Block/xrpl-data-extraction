package xrpl

type XRPLManager struct {
	HTTPClient *HTTPClient
	WSClient   *WebSocketClient
}

// NewXRPLManager creates a new XRPL manager with the provided base URL and WebSocket URL
func NewXRPLManager(baseURL, wsURL string) (*XRPLManager, error) {
	// Initialize the HTTP AND WS client
	httpClient := NewHTTPClient(baseURL)
	wsClient, err := NewWebSocketClient(wsURL)
	if err != nil {
		return nil, err
	}
	// Return the XRPL manager with the HTTP and WS clients combined
	return &XRPLManager{
		HTTPClient: httpClient,
		WSClient:   wsClient,
	}, nil
}

// GetHTTPClient returns the HTTP client
func (x *XRPLManager) GetHTTPClient() *HTTPClient {
	return x.HTTPClient
}

// GetWSClient returns the WebSocket client
func (x *XRPLManager) GetWSClient() *WebSocketClient {
	return x.WSClient
}
