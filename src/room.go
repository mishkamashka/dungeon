package main

type Room struct {
	id int
	name string
	clients []*Client
	description string
}
