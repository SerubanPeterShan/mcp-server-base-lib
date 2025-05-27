package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	config := &Config{
		Port:   "8080",
		Logger: logrus.New(),
		Handlers: map[string]MessageHandler{
			"test": func(s *Server, conn *websocket.Conn, msg Message) error {
				return nil
			},
		},
	}

	server := NewServer(config)
	assert.NotNil(t, server)
	assert.NotNil(t, server.clients)
	assert.NotNil(t, server.broadcast)
	assert.NotNil(t, server.register)
	assert.NotNil(t, server.unregister)
	assert.NotNil(t, server.handlers)
}

func TestServer_RegisterHandler(t *testing.T) {
	server := NewServer(&Config{Logger: logrus.New()})

	handler := func(s *Server, conn *websocket.Conn, msg Message) error {
		return nil
	}

	server.RegisterHandler("test", handler)
	assert.NotNil(t, server.handlers["test"])
}

func TestServer_Broadcast(t *testing.T) {
	server := NewServer(&Config{Logger: logrus.New()})

	// Create a test message
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "test"}`),
	}

	// Start the server
	go server.handleWebSocket()

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(server.handleWebSocketConnection))
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Register the client
	server.register <- conn

	// Wait for client to be registered
	time.Sleep(100 * time.Millisecond)

	// Broadcast the message
	server.Broadcast(msg)

	// Read the message from the client
	var receivedMsg Message
	err = conn.ReadJSON(&receivedMsg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Type, receivedMsg.Type)
	assert.Equal(t, msg.Payload, receivedMsg.Payload)
}

func TestServer_Send(t *testing.T) {
	server := NewServer(&Config{Logger: logrus.New()})

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(server.handleWebSocketConnection))
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Create a test message
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "test"}`),
	}

	// Send the message
	err = server.Send(conn, msg)
	assert.NoError(t, err)

	// Read the message from the client
	var receivedMsg Message
	err = conn.ReadJSON(&receivedMsg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Type, receivedMsg.Type)
	assert.Equal(t, msg.Payload, receivedMsg.Payload)
}

func TestServer_GetClients(t *testing.T) {
	server := NewServer(&Config{Logger: logrus.New()})

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(server.handleWebSocketConnection))
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Register the client
	server.register <- conn

	// Wait for client to be registered
	time.Sleep(100 * time.Millisecond)

	// Get clients
	clients := server.GetClients()
	assert.Len(t, clients, 1)
}

func TestServer_HealthCheck(t *testing.T) {
	server := NewServer(&Config{Logger: logrus.New()})

	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Call the health check handler
	server.handleHealthCheck(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}
