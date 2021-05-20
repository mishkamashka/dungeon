package main

func main() {

	server := NewServer()
	server.rooms[0].description = "This is a Hall - a wide room with columns. As you look around you see a fountain in a middle of the room. You also notice a pass to another room with is located North of you."
	server.rooms[1].description = "You're standing in front of the turnstile, on your right you can see a small cabin and a man inside. He is staring at a screen in front of him and doesn't pay any attention to you. It seems that you need a special pass-card to go further."
	server.Listen("localhost:8000")

}
