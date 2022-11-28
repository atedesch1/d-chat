package chat

import (
	"context"
	"errors"
	"log"
	"net"

	chat_message "github.com/decentralized-chat/pb"
	"github.com/decentralized-chat/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) RegisterServer() {
	var err error
	c.lis, err = net.Listen("tcp", util.JoinIpAndPort(c.User.Addr.Ip, int(c.User.Addr.Port)))

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
	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, c.DialUser(msg.User)
}

func (c *Client) DialUser(user *chat_message.User) error {
	isAlreadyConnected := false
	for _, peer := range c.peers {
		if peer.user.Addr.Ip == user.Addr.Ip &&
			peer.user.Addr.Port == user.Addr.Port {
			isAlreadyConnected = true
			break
		}
	}

	if isAlreadyConnected {
		return errors.New("user already connected")
	}

	conn, err := grpc.Dial(
		util.JoinIpAndPort(user.Addr.Ip, int(user.Addr.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}

	c.peersMutex.Lock()
	c.peers = append(c.peers, Peer{
		user:   user,
		conn:   conn,
		client: chat_message.NewChatServiceClient(conn),
	})
	c.peersMutex.Unlock()

	return nil
}

func (c *Client) DialChannel(channel Channel) {
	for _, user := range channel.users {
		go c.DialUser(user)
	}
}
