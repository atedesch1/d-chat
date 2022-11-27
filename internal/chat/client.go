package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (c *Client) SendMessage(ctx context.Context, msg *chat_message.ContentMessage) (*chat_message.AckMessage, error) {
	fmt.Println(msg.From.Username, ":", msg.Content)

	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, nil
}

func (c *Client) BroadcastMessage(content string) error {
	msg := &chat_message.ContentMessage{
		From:    &c.User,
		Content: content,
		SentAt:  timestamppb.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, peer := range c.peers {
		ack, err := peer.client.SendMessage(ctx, msg)
		if err != nil {
			return err
		}
		log.Println("Ack from", ack.From.Username)
	}

	return nil
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
