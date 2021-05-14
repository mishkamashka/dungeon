package main

func main() {

	server := NewServer()
	server.Listen("localhost:8000")

}
