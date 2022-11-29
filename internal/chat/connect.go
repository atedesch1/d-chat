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
	"google.golang.org/protobuf/types/known/emptypb"
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
	// log.Printf("listening on %v", c.lis.Addr())

	if err := c.srv.Serve(c.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (c *Client) DialAddress(addr *chat_message.Address) error {
	for _, peer := range c.peers {
		if peer.user.Addr.Ip == addr.Ip && peer.user.Addr.Port == addr.Port {
			return errors.New("user already connected")
		}
	}

	conn, err := grpc.Dial(
		util.JoinIpAndPort(addr.Ip, int(addr.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}

	client := chat_message.NewChatServiceClient(conn)

	userInfo, err := client.GetUsername(context.Background(), &emptypb.Empty{})

	if err != nil {
		return err
	}

	c.AddPeer(userInfo.User, conn)

	_, err = client.RequestConnection(context.Background(), &chat_message.ConnectionMessage{
		User: &c.User,
	})

	return err
}

func (c *Client) RequestMatchConnection(username string) error {
	peer, ok := c.peers[username]

	if !ok {
		return errors.New("peer not connected")
	}

	if _, err := peer.client.RequestConnection(context.Background(), &chat_message.ConnectionMessage{
		User: &c.User,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Client) CloseConnection(username string) error {
	peer, ok := c.peers[username]
	if !ok {
		return errors.New("peer not connected")
	}

	if _, err := peer.client.Disconnect(context.Background(), &chat_message.ConnectionMessage{
		User: &c.User,
	}); err != nil {
		return err
	}

	peer.conn.Close()
	return c.RemovePeer(username)
}
