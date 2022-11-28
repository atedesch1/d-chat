package main

import (
	"flag"

	"github.com/decentralized-chat/internal/chat"
)

var (
	username = flag.String("username", "user", "The username")
	port     = flag.Uint("port", 50000, "The server port")
)

func main() {
	flag.Parse()

	c := chat.NewClient(*username, *port)

	c.RegisterServer()

	go c.ListenForConnections()

	c.ListenForInput()
}
