package main

import (
	"net"
	"fmt"
	"os"
	"github.com/labstack/gommon/log"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	log.Println("Connection Handle")

	buf := make([]byte, 128)

	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	log.Println("Rec ", string(buf))

	if (string(buf) == "SET") {
		conn.Write([]byte("+PONG"))
	}

	// Close the connection when you're done with it.
	//conn.Close()
}
