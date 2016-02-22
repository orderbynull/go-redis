package main

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"github.com/labstack/gommon/log"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type Client struct {
	outgoing chan string
	writer   *bufio.Writer
	connection *net.Conn
}

func (client *Client) Write() {
	for data := range client.outgoing {
		log.Print("Write ", data)
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go func() {
		reader := bufio.NewReader(*client.connection)
		reader.ReadLine()

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Print("Cannot read error", err)

				connection := *client.connection
				connection.Close()
				return
			}

			log.Print("Read ", strings.Trim(line, "\r\n"))

			if (line == "PING") {
				client.outgoing <- "+PONG\r\n"
			} else if (line == "PONG") {
				client.outgoing <- "+PING\r\n"
			} else if (line == "SET") {
				value, errValue := reader.ReadString('\n')
				if (errValue != nil) {
					log.Print("Cannot read error", err)

					connection := *client.connection
					connection.Close()
					return
				}

				client.outgoing <- value
				client.outgoing <- "+OK\r\n"
			} else {
				client.outgoing <- "+OK\r\n"
			}
		}
	}()

	go client.Write()
}

func NewClient(connection *net.Conn) *Client {
	writer := bufio.NewWriter(*connection)

	client := &Client{
		outgoing: make(chan string, 1),
		writer: writer,
		connection: connection,
	}

	client.Listen()

	return client
}

type Clients struct {
	clients []*Client
}

func (clients *Clients) Join(connection *net.Conn) {
	client := NewClient(connection)
	clients.clients = append(clients.clients, client)
}

var (
	ClientsStack *Clients;
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()


	ClientsStack = new(Clients)

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}


		ClientsStack.Join(&conn)
	}
}