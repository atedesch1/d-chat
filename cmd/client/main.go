package main

import (
	"flag"
	"fmt"

	tm "github.com/buger/goterm"
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

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()

	fmt.Println("Decentralized-Chat")

	go c.ListenForConnections()

	c.ListenForInput()
}
