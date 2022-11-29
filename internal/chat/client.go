package chat

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/decentralized-chat/internal/server"
	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
)

type Client struct {
	chat_message.UnimplementedChatServiceServer

	User chat_message.User

	lis net.Listener
	srv *grpc.Server

	peersMutex sync.Mutex
	peers      map[string]Peer
	channel    string

	zk *server.Server
}

func NewClient(username string, ip string, port uint) *Client {
	zk := new(server.Server)
	zk.Init(ip, "2181")

	return &Client{
		User: chat_message.User{
			Username: username,
			Addr: &chat_message.Address{
				Ip:   ip,
				Port: uint32(port),
			},
		},
		peers:   make(map[string]Peer),
		channel: "",
		zk:      zk,
	}
}

type Peer struct {
	user   *chat_message.User
	client chat_message.ChatServiceClient
	conn   *grpc.ClientConn
}

func (c *Client) AddPeer(user *chat_message.User, conn *grpc.ClientConn) {
	c.peersMutex.Lock()
	c.peers[user.Username] = Peer{
		user:   user,
		conn:   conn,
		client: chat_message.NewChatServiceClient(conn),
	}
	c.peersMutex.Unlock()
}

func (c *Client) RemovePeer(username string) error {
	if _, ok := c.peers[username]; !ok {
		return errors.New("couldnt find peer")
	}

	c.peersMutex.Lock()
	delete(c.peers, username)
	c.peersMutex.Unlock()

	return nil
}

func (c *Client) ListenForInput() {
	inputChannel := make(chan string)

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

		if strings.HasPrefix(input, "$") {
			_, in, _ := strings.Cut(input, "$")
			command, params, _ := strings.Cut(in, " ")

			switch command {
			case "create":
				c.CreateChannel(params)
			case "list":
				fmt.Println("Channels:", c.ListChannels())
			case "join":
				c.JoinChannel(params)
			case "leave":
				c.DisconnectFromChannel()
			}

		} else {
			go c.BroadcastMessage(input)
		}
	}
}
