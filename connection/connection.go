package connection

import (
	"fmt"
	"net"

	"github.com/pbharrell/bloner-server/connection"
)

func connectToServer() connection.Server {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err == nil {
		fmt.Println("Error connecting to server:", err)
	}

	s := connection.Server{
		Conn: conn,
		Data: make(chan connection.Message),
	}

	go s.Listen()

	return s
}
