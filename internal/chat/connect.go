package chat

import (
	"fmt"
	"log"
	"net"
	"strconv"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	log.Printf("server listening at %v", c.lis.Addr())
	if err := c.srv.Serve(c.lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (c *Client) DialChannel(channel Channel) {
	for _, user := range channel.users {
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(user.Addr.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c.peers = append(c.peers, Peer{
			user:   user,
			conn:   conn,
			client: chat_message.NewChatServiceClient(conn),
		})
	}
}
