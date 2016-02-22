package main

import (
	"net"
	"fmt"
	"os"
	"github.com/labstack/gommon/log"
"bufio"
"bytes"
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


	//readbuffer := bytes.NewBuffer([]byte("123\r\n456"))
	//reader := bufio.NewReader(readbuffer)

	bufferForReader := bytes.NewBuffer(buf)
	bufReader := bufio.NewReader(bufferForReader)

	read:
	line, isPrefis, readErr := bufReader.ReadLine();
	for  readErr == nil {
		log.Println("Request line ", string(line))
		log.Println(isPrefis)

		//log.Print("PING OPERATROR %s", string(line) == "PING")
		if (string(line) == "PING") {
			writer := bufio.NewWriter(conn)
			writer.WriteString("+PONG\r\n")
			writer.Flush()
		} else if (string(line) == "PONG") {
			writer := bufio.NewWriter(conn)
			writer.WriteString("+PING\r\n")
			writer.Flush()
		} else {
			writer := bufio.NewWriter(conn)
			writer.WriteString("+OK\r\n")
			writer.Flush()
		}

		goto read
	}

	//numbs, responseWriteError := conn.Write([]byte("+PONG"))
	//log.Print(numbs, responseWriteError)


	// Close the connection when you're done with it.
	//conn.Close()
}
