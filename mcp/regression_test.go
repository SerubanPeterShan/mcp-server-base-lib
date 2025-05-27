package mcp

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegressionReconnection(t *testing.T) {
	// Test client reconnection behavior
	server := NewServer(&Config{Logger: logrus.New()})
	go server.Start("8087")
	waitForServerReady("8087", t)

	// First connection
	conn1, _, err := websocket.DefaultDialer.Dial("ws://localhost:8087/ws", nil)
	require.NoError(t, err)
	conn1.Close()

	// Second connection
	conn2, _, err := websocket.DefaultDialer.Dial("ws://localhost:8087/ws", nil)
	require.NoError(t, err)
	defer conn2.Close()

	// Verify client list
	clients := server.GetClients()
	assert.Len(t, clients, 1)
}

func TestRegressionMessageOrder(t *testing.T) {
	// Test message ordering and delivery
	server := NewServer(&Config{Logger: logrus.New()})
	go server.Start("8088")
	waitForServerReady("8088", t)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8088/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send multiple messages
	messages := []Message{
		{Type: "test", Payload: json.RawMessage(`{"order": 1}`)},
		{Type: "test", Payload: json.RawMessage(`{"order": 2}`)},
		{Type: "test", Payload: json.RawMessage(`{"order": 3}`)},
	}

	for _, msg := range messages {
		err := server.Send(conn, msg)
		require.NoError(t, err)
	}

	// Verify messages are received in order
	for i := 0; i < len(messages); i++ {
		var receivedMsg Message
		err := conn.ReadJSON(&receivedMsg)
		require.NoError(t, err)
		assert.Equal(t, messages[i].Type, receivedMsg.Type)
		assert.Equal(t, messages[i].Payload, receivedMsg.Payload)
	}
}

func TestRegressionErrorHandling(t *testing.T) {
	// Test error handling scenarios
	server := NewServer(&Config{
		Logger: logrus.New(),
		Handlers: map[string]MessageHandler{
			"error": func(s *Server, conn *websocket.Conn, msg Message) error {
				return errors.New("simulated error for testing purposes")
			},
		},
	})
	go server.Start("8089")
	waitForServerReady("8089", t)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8089/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send message that triggers error
	msg := Message{
		Type:    "error",
		Payload: json.RawMessage(`{"data": "error test"}`),
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)

	// Verify connection is still alive
	err = conn.WriteJSON(Message{Type: "test", Payload: json.RawMessage(`{"data": "test"}`)})
	require.NoError(t, err)
}

func TestRegressionConcurrentConnections(t *testing.T) {
	// Test handling of many concurrent connections
	server := NewServer(&Config{Logger: logrus.New()})
	go server.Start("8090")
	waitForServerReady("8090", t)

	// Create many connections
	connections := make([]*websocket.Conn, 50)
	for i := range connections {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8090/ws", nil)
		require.NoError(t, err)
		connections[i] = conn
		defer conn.Close()
	}

	// Verify all connections are registered
	clients := server.GetClients()
	assert.Len(t, clients, len(connections))

	// Test broadcasting to all connections
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "broadcast test"}`),
	}
	server.Broadcast(msg)

	// Verify all connections receive the message
	for _, conn := range connections {
		var receivedMsg Message
		err := conn.ReadJSON(&receivedMsg)
		require.NoError(t, err)
		assert.Equal(t, msg.Type, receivedMsg.Type)
		assert.Equal(t, msg.Payload, receivedMsg.Payload)
	}
}
