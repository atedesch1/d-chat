package main

import (
	"flag"

	"github.com/decentralized-chat/internal/chat"
)

var (
	username = flag.String("username", "user", "The username")
	ip       = "127.0.0.1"
	port     = flag.Uint("port", 50000, "The server port")
)

func main() {
	flag.Parse()

	c := chat.NewClient(*username, ip, *port)

	c.RegisterUser()

	c.RegisterServer()

	go c.ListenForConnections()

	c.ListenForInput()
}
