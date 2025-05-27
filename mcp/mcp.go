package mcp

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Message represents a message in the MCP protocol
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Server represents an MCP server instance
type Server struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
	logger     *logrus.Logger
	handlers   map[string]MessageHandler
	server     *http.Server
}

// MessageHandler is a function type for handling messages
type MessageHandler func(*Server, *websocket.Conn, Message) error

// Config holds the server configuration
type Config struct {
	Port     string
	Logger   *logrus.Logger
	Handlers map[string]MessageHandler
}

// NewServer creates a new MCP server instance
func NewServer(config *Config) *Server {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		logger:     config.Logger,
		handlers:   config.Handlers,
	}
}

// RegisterHandler registers a new message handler
func (s *Server) RegisterHandler(messageType string, handler MessageHandler) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers[messageType] = handler
}

// Broadcast sends a message to all connected clients
func (s *Server) Broadcast(msg Message) {
	s.broadcast <- msg
}

// Send sends a message to a specific client
func (s *Server) Send(client *websocket.Conn, msg Message) error {
	return client.WriteJSON(msg)
}

// GetClients returns a list of all connected clients
func (s *Server) GetClients() []*websocket.Conn {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clients := make([]*websocket.Conn, 0, len(s.clients))
	for client := range s.clients {
		clients = append(clients, client)
	}
	return clients
}

// Start initializes and starts the MCP server
func (s *Server) Start(port string) error {
	// Start WebSocket handler
	go s.handleWebSocket()

	// Set up HTTP routes
	http.HandleFunc("/ws", s.handleWebSocketConnection)
	http.HandleFunc("/health", s.handleHealthCheck)

	// Start HTTP server
	s.logger.Infof("Starting MCP server on port %s", port)
	s.server = &http.Server{Addr: ":" + port}
	return s.server.ListenAndServe()
}

// handleWebSocketConnection handles new WebSocket connections
func (s *Server) handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	s.register <- conn

	// Handle incoming messages
	go func() {
		defer func() {
			s.unregister <- conn
			conn.Close()
		}()

		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				s.logger.Errorf("Error reading message: %v", err)
				break
			}

			// Handle message if a handler exists
			if handler, ok := s.handlers[msg.Type]; ok {
				if err := handler(s, conn, msg); err != nil {
					s.logger.Errorf("Error handling message: %v", err)
				}
			} else {
				// If no handler exists, broadcast the message
				s.broadcast <- msg
			}
		}
	}()
}

// handleWebSocket manages WebSocket connections and message broadcasting
func (s *Server) handleWebSocket() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			s.logger.Infof("New client connected")

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				s.logger.Infof("Client disconnected")
			}
			s.mutex.Unlock()

		case message := <-s.broadcast:
			s.mutex.Lock()
			for client := range s.clients {
				err := client.WriteJSON(message)
				if err != nil {
					s.logger.Errorf("Error broadcasting message: %v", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mutex.Unlock()
		}
	}
}

// handleHealthCheck provides a health check endpoint
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Close all client connections
	for client := range s.clients {
		client.Close()
		delete(s.clients, client)
	}

	// Close channels
	close(s.broadcast)
	close(s.register)
	close(s.unregister)

	return nil
}
