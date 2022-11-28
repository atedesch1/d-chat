package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
)

type Client struct {
	chat_message.UnimplementedChatServiceServer

	User chat_message.User

	lis net.Listener
	srv *grpc.Server

	peers []Peer
}

func NewClient(username string, port uint) *Client {
	return &Client{
		User: chat_message.User{
			Username: username,
			Addr: &chat_message.Address{
				Ip:   "localhost",
				Port: uint32(port),
			},
		},
	}
}

func (c *Client) ListenForInput() {
	inputChannel := make(chan string)

	// Read input
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _, err := reader.ReadLine()
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
			inputChannel <- string(input)
		}
	}()

	for input := range inputChannel {
		go c.BroadcastMessage(input)
	}
}
