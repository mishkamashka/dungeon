package main

import (
	"io"
	"net"
)

type Client struct {
	connection   net.Conn
	name   string
	writer io.Writer
}
