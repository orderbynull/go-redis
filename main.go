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

type StringMap map[string]string

type DataBase struct {
	Strings StringMap
}

var (
	DataBases [10]DataBase
)

type Client struct {
	outgoing   chan string
	writer     *bufio.Writer
	connection *net.Conn
	DataBase   DataBase
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

		client.outgoing <- "+OK\r\n"

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Print("Cannot read error", err)

				connection := *client.connection
				connection.Close()
				return
			}

			line = strings.Trim(line, "\r\n")
			log.Println(fmt.Sprintf("Read %s", line))

			if (line == "PING") {
				client.outgoing <- "+PONG\r\n"
			} else if (line == "PONG") {
				client.outgoing <- "+PING\r\n"
			} else if (line == "GET") {
				reader.ReadString('\n')
				key, errValue := reader.ReadString('\n')
				if (errValue != nil) {
					log.Print("Cannot read error", err)

					connection := *client.connection
					connection.Close()
					return
				}
				key = strings.Trim(key, "\r\n")

				value := client.DataBase.Strings[key]

				returnValue := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				client.outgoing <- returnValue
				//client.outgoing <- "+OK\r\n"
			} else if (line == "SET") {
				reader.ReadString('\n')

				key, errValue := reader.ReadString('\n')
				if (errValue != nil) {
					log.Print("Cannot read error", err)

					connection := *client.connection
					connection.Close()
					return
				}
				key = strings.Trim(key, "\r\n")
				log.Println(fmt.Sprintf("Read key '%s'", key))

				reader.ReadString('\n')

				value, errValue := reader.ReadString('\n')
				if (errValue != nil) {
					log.Print("Cannot read error", err)

					connection := *client.connection
					connection.Close()
					return
				}
				value = strings.Trim(value, "\r\n")
				log.Println(fmt.Sprintf("Read value '%s'", value))

				client.DataBase.Strings[key] = value
				//client.outgoing <- "+OK\r\n"
			} else {

			}
		}
	}()

	go client.Write()
}

func NewClient(connection *net.Conn, db DataBase) *Client {
	writer := bufio.NewWriter(*connection)

	client := &Client{
		outgoing: make(chan string, 100),
		writer: writer,
		connection: connection,
		DataBase: db,
	}

	client.Listen()

	return client
}

type Clients struct {
	clients []*Client
}

func (clients *Clients) Join(connection *net.Conn) {
	client := NewClient(connection, DataBases[0])
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

	DataBases = [10]DataBase{
		DataBase{
			Strings: StringMap{},
		}, // 1
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
		DataBase{
			Strings: StringMap{},
		},
	}

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