package main

import "strconv"

type Room struct {
	id int
	name string
	clients []*Client
	description string

	//todo actual items logic here
	items string
	actions string
	//isFight
}

func (room *Room) describe() string {
	fullDescription := room.description
	people := "\n\nYou can see " + strconv.Itoa(len(room.clients)) + " persons in here. Seems that nothing grand is happening." //isFight check
	fullDescription = fullDescription + people + room.items + room.actions
	return fullDescription
}