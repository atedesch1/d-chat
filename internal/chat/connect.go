package chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) RegisterServer() {
	var err error
	c.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", c.User.Addr.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	c.srv = grpc.NewServer()
	chat_message.RegisterChatServiceServer(c.srv, c)
}

func (c *Client) ListenForConnections() {
	log.Printf("listening on %v", c.lis.Addr())
	if err := c.srv.Serve(c.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (c *Client) RequestConnection(ctx context.Context, msg *chat_message.ConnectionMessage) (*chat_message.AckMessage, error) {
	isAlreadyConnected := false
	for _, peer := range c.peers {
		if peer.user.Addr.Ip == msg.ConnectTo.Addr.Ip && peer.user.Addr.Port == msg.ConnectTo.Addr.Port {
			isAlreadyConnected = true
			break
		}
	}

	var err error

	if !isAlreadyConnected {
		c.DialUser(msg.ConnectTo)
	} else {
		log.Println("is already connected")
		err = errors.New("client is already connected")
	}

	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, err
}

func (c *Client) DialUser(user *chat_message.User) error {
	conn, err := grpc.Dial(
		"localhost:"+strconv.Itoa(int(user.Addr.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}

	c.peers = append(c.peers, Peer{
		user:   user,
		conn:   conn,
		client: chat_message.NewChatServiceClient(conn),
	})

	return nil
}

func (c *Client) DialChannel(channel Channel) {
	for _, user := range channel.users {
		c.DialUser(user)
	}
}
