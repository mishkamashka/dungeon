package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"regexp"
	"sync"
)

type Server struct {
	listener net.Listener
	rooms    [2]*Room
	mutex    *sync.Mutex
}

func NewServer() *Server {
	return &Server{
		rooms: [2]*Room{{id: 0, name: "Hall", clients: []*Client{}}, {id: 1, name: "Security", clients: []*Client{}}},
		mutex: &sync.Mutex{},
	}
}

func (server *Server) Start(address string) {
	server.Listen(address)
}

func (server *Server) Listen(address string) {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		log.Println(err)
		return
	}

	server.listener = ln
	log.Printf("Listening on %v", address)

	for {
		conn, err := server.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			client := server.accept(conn)
			go server.serve(client)
		}
	}
}

func (server *Server) accept(conn net.Conn) *Client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(server.rooms[0].clients)+len(server.rooms[1].clients)+1)

	client := &Client{
		connection: conn,
		writer:     io.Writer(conn),
	}

	client.name = getClientNameAndClearTerminal(client)
	client.room = 0

	server.mutex.Lock()
	defer server.mutex.Unlock()
	server.rooms[client.room].clients = append(server.rooms[client.room].clients, client)

	server.send(client, "welcome to room " + server.rooms[client.room].name)
	server.send(client, server.rooms[client.room].description)
	server.broadcastInRoomExceptSender(client.name, client.room, client.name+" has joined")
	showPrompt(client)
	return client
}

func getClientNameAndClearTerminal(client *Client) string {
	client.writer.Write([]byte("\033[1;1H\033[2J"))
	_, err := client.writer.Write([]byte("Enter your name:\n"))
	if err != nil {
		log.Println(err)
	}

	name := ""
	for name == "" || name == "\n" {
		name, err = bufio.NewReader(client.connection).ReadString('\n')
		re := regexp.MustCompile(`\r?\n`)
		name = re.ReplaceAllString(name, "")

		if err != nil {
			err := client.connection.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}

	return name
}

func (server *Server) remove(client *Client) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for i, item := range server.rooms[client.room].clients {
		if item == client {
			server.rooms[client.room].clients = append(server.rooms[client.room].clients[:i], server.rooms[client.room].clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.connection.RemoteAddr().String())
	server.broadcast(client.name + " has left the server")
	client.connection.Close()
}

func (server *Server) serve(client *Client) {
	defer server.remove(client)

	for {
		message, err := bufio.NewReader(client.connection).ReadString('\n')
		re := regexp.MustCompile(`\r?\n`)
		message = re.ReplaceAllString(message, "")

		// todo command logic

		if err != nil {
			err := client.connection.Close()
			if err != nil {
				log.Println(err)
			}
			return
		}

		server.runCommand(message, client)

		showPrompt(client)
	}
}


func (server *Server) send(client *Client, message string) {
	client.writer.Write([]byte("\r" + message + "\n"))
}

func (server *Server) broadcastInRoom(room int, message string) {
	for _, client := range server.rooms[room].clients {
		client.writer.Write([]byte("\r" + message + "\n> "))
	}
	return
}

func (server *Server) broadcast(message string) {
	for _, room := range server.rooms {
		for _, client := range room.clients {
			client.writer.Write([]byte("\r" + message + "\n> "))
		}
	}
	return
}

func (server *Server) broadcastInRoomExceptSender(name string, room int, message string) {
	for _, client := range server.rooms[room].clients {
		if client.name != name {
			client.writer.Write([]byte("\r" + message + "\n> "))
		}
	}
	return
}

func showPrompt(client *Client) {
	client.writer.Write([]byte("> "))
}

func (server *Server) runCommand(command string, client *Client) {
	switch command {
	case "N", "S":
		server.switchRoom(command, client)
	default:
		if command != "" {
			server.broadcastInRoomExceptSender(client.name, client.room, client.name+" says: "+command)
		}
	}
}

func (server *Server) switchRoom(direction string, client *Client) {
	var newRoom int
	if direction == "N" && client.room == 0 {
		newRoom = 1
	} else if direction == "S" && client.room == 1 {
		newRoom = 0
	} else {
		server.send(client, "You cannot move there, other room in an another direction")
		return
	}


	prevRoom := client.room
	client.room = newRoom

	server.broadcastInRoomExceptSender(client.name, prevRoom, client.name + " has left the room")
	server.broadcastInRoomExceptSender(client.name, newRoom, client.name + " has entered the room")
	server.send(client,  "You've entered room " + server.rooms[newRoom].name)
	server.send(client, server.rooms[newRoom].description)

	server.mutex.Lock()
	defer server.mutex.Unlock()

	for i, item := range server.rooms[prevRoom].clients {
		if item == client {
			server.rooms[prevRoom].clients = append(server.rooms[prevRoom].clients[:i], server.rooms[prevRoom].clients[i+1:]...)
		}
	}
	server.rooms[newRoom].clients = append(server.rooms[newRoom].clients, client)
}
