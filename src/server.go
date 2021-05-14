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
	clients  []*Client
	mutex    *sync.Mutex
}

func NewServer() *Server {
	return &Server{
		mutex: &sync.Mutex{},
	}
}

func (server *Server) Start(address string) {
	server.Listen(address)
}

func (server *Server) Listen(address string){
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
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(server.clients)+1)

	client := &Client{
		connection: conn,
		writer: io.Writer(conn),
	}

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

	client.name = name
	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.clients = append(server.clients, client)

	server.send(client.name, "welcome")
	server.broadcastExceptSender(client.name, client.name + " has joined")
	return client
}

func (server *Server) remove(client *Client) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for i, item := range server.clients {
		if item == client {
			server.clients = append(server.clients[:i], server.clients[i+1:] ...)
		}
	}

	log.Printf("Closing connection from %v", client.connection.RemoteAddr().String())
	server.broadcastExceptSender(client.name, client.name + " has left")
	client.connection.Close()
}


func (server *Server) serve(client *Client) {
	defer server.remove(client)

	for {
		message, err := bufio.NewReader(client.connection).ReadString('\n')
		re := regexp.MustCompile(`\r?\n`)
		message = re.ReplaceAllString(message, "")

		if err != nil {
			err := client.connection.Close()
			if err != nil {
				log.Println(err)
			}
			return
		}
		server.broadcastExceptSender(client.name, client.name + " says: " + message)
	}
}

func (server *Server) broadcast(message string) {
	for _, client := range server.clients {
		// TODO: handle error here?
		client.writer.Write([]byte(message + "\n"))
	}
	return
}

func (server *Server) broadcastExceptSender(name, message string) {
	for _, client := range server.clients {
		// TODO: handle error here?
		if client.name != name {
			client.writer.Write([]byte(message + "\n"))
		}
	}
	return
}

func (server *Server) send(name, message string) {
	for _, client := range server.clients {
		if client.name == name {
			client.writer.Write([]byte(message + "\n"))
			break
		}
	}
	return
}

