package mcp

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmokeBasicFunctionality(t *testing.T) {
	// Basic smoke test to verify core functionality
	server := NewServer(&Config{Logger: logrus.New()})
	defer server.Stop()

	// Start server
	go func() {
		err := server.Start("0") // Use dynamic port allocation
		require.NoError(t, err)
	}()
	port := server.GetPort() // Retrieve the dynamically assigned port
	waitForServerReady(port, t)

	// Test WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:"+port+"/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Test message sending
	msg := Message{
		Type:    "test",
		Payload: json.RawMessage(`{"data": "smoke test"}`),
	}
	err = server.Send(conn, msg)
	require.NoError(t, err)

	// Verify message received
	var receivedMsg Message
	err = conn.ReadJSON(&receivedMsg)
	require.NoError(t, err)
	assert.Equal(t, msg.Type, receivedMsg.Type)
	assert.Equal(t, msg.Payload, receivedMsg.Payload)
}

func TestSmokeServerStartup(t *testing.T) {
	// Test server startup with different configurations
	testCases := []struct {
		name   string
		config *Config
		port   string
	}{
		{
			name: "Default Config",
			config: &Config{
				Logger: logrus.New(),
			},
			port: "8084",
		},
		{
			name: "With Handlers",
			config: &Config{
				Logger: logrus.New(),
				Handlers: map[string]MessageHandler{
					"echo": func(s *Server, conn *websocket.Conn, msg Message) error {
						return s.Send(conn, msg)
					},
				},
			},
			port: "8085",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(tc.config)
			go func() {
				err := server.Start(tc.port)
				require.NoError(t, err)
			}()
			waitForServerReady(tc.port, t)
		})
	}
}

func TestSmokeMessageHandling(t *testing.T) {
	// Test basic message handling functionality
	server := NewServer(&Config{
		Logger: logrus.New(),
		Handlers: map[string]MessageHandler{
			"echo": func(s *Server, conn *websocket.Conn, msg Message) error {
				return s.Send(conn, msg)
			},
		},
	})

	go server.Start("8086")
	waitForServerReady("8086", t)

	// Connect client
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8086/ws", nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send message
	msg := Message{
		Type:    "echo",
		Payload: json.RawMessage(`{"data": "echo test"}`),
	}
	err = conn.WriteJSON(msg)
	require.NoError(t, err)

	// Verify echo response
	var response Message
	err = conn.ReadJSON(&response)
	require.NoError(t, err)
	assert.Equal(t, msg.Type, response.Type)
	assert.Equal(t, msg.Payload, response.Payload)
}
