package main

import (
	"flag"

	"github.com/decentralized-chat/internal/chat"
	chat_message "github.com/decentralized-chat/pb"
)

var (
	username = flag.String("username", "default", "The username")
	cport    = flag.Uint("cport", 50001, "The client port")
	sport    = flag.Uint("sport", 50002, "The server port")
)

func main() {
	flag.Parse()

	c := chat.NewClient(*username, *cport)

	channel := chat.NewChannel(0)
	channel.AddUser(chat_message.User{
		Username: "default",
		Addr: &chat_message.Address{
			Ip:   "localhost",
			Port: uint32(*sport),
		},
	})

	c.RegisterServer()

	go c.ListenForConnections()
	go c.DialChannel(*channel)

	c.ListenForInput()
}
