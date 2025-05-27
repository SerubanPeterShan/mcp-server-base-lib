package mcp

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const skipIntegrationTestMsg = "Skipping integration test in short mode"

func TestIntegrationMultipleClients(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip(skipIntegrationTestMsg)
	}

	server := NewServer(&Config{
		Logger: logrus.New(),
		Handlers: map[string]MessageHandler{
			"echo": func(s *Server, conn *websocket.Conn, msg Message) error {
				return s.Send(conn, msg)
			},
		},
	})

	// Start server
	go func() {
		err := server.Start("8080")
		require.NoError(t, err)
	}()

	// Connect multiple clients
	clients := make([]*websocket.Conn, 3)
	for i := range clients {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		require.NoError(t, err)
		clients[i] = conn
		defer conn.Close()
	}

	// Test message broadcasting
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "broadcast"}`),
	}

	server.Broadcast(msg)

	// Verify all clients receive the message
	for _, client := range clients {
		var receivedMsg Message
		err := client.ReadJSON(&receivedMsg)
		require.NoError(t, err)
		assert.Equal(t, msg.Type, receivedMsg.Type)
		assert.Equal(t, msg.Payload, receivedMsg.Payload)
	}
}

func TestIntegrationConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTestMsg)
	}

	server := NewServer(&Config{Logger: logrus.New()})
	go server.Start("8081")

	// Create multiple clients
	clients := make([]*websocket.Conn, 5)
	for i := range clients {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
		require.NoError(t, err)
		clients[i] = conn
		defer conn.Close()
	}

	// Test concurrent message sending
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			msg := Message{
				Type:    "test",
				Payload: json.RawMessage(`{"index": ` + strconv.Itoa(index) + `}`),
			}
			server.Broadcast(msg)
			done <- true
		}(i)
	}

	// Wait for all messages to be sent
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestIntegrationServerShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip(skipIntegrationTestMsg)
	}

	server := NewServer(&Config{Logger: logrus.New()})
	go server.Start("8082")

	// Connect a client
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8082/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Verify client is registered
	clients := server.GetClients()
	assert.Len(t, clients, 1)

	// Properly shut down the server
	err = server.Stop()
	require.NoError(t, err)

	// Verify client is disconnected
	clients = server.GetClients()
	assert.Len(t, clients, 0)
}
