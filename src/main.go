package main

func main() {

	server := NewServer()

	//move it somewhere ltr
	server.rooms[0].description = "This is a Hall - a wide room with columns. As you look around you see a fountain in a middle of the room. You also notice a pass to another room which is located North of you."
	server.rooms[0].items = "\n\nThere are several information stands, probably some events are announced there. You can see a bunch of benches around the columns."
	server.rooms[0].actions = "\n\nSeems that you can walk towards the security cabin, look around this place or put your finger in a fountain."

	server.rooms[1].description = "You're standing in front of the turnstile, on your right you can see a small cabin and a man inside. He is staring at a screen in front of him and doesn't pay any attention to you. It seems that you need a special pass-card to go further."
	server.Listen("localhost:8000")

}
