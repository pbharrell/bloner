package connection

import (
	"fmt"
	"net"
)

type Message struct {
	// Suppported types:
	// Lobby Types:
	//   MO: lobby_req; data = lobby_id
	//   MT: lobby_assign; data = { lobby_id, player_id }
	//
	// Game Init Types:
	//   MT: game_start; data = player_id
	//
	// State Types:
	//   MT: state_req; data = nil
	//   MO: state_res; data = gameState (full)
	//   MT: state_res; data = gameState (full)
	//
	// Turn Types:
	//   MO: state_update; data = gameState (changed)
	//   MT: state_update; data = gameState (changed)
	//

	Type string `json:"type"`
	Data any    `json:"data"` // payload
}

type Server struct {
	Conn net.Conn
	Data chan Message
}

func connectToServer() Server {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err == nil {
		fmt.Println("Error connecting to server:", err)
	}

	s := Server{
		Conn: conn,
		Data: make(chan Message),
	}

	go s.listen()

	return s
}

func (s *Server) listen() {

}
