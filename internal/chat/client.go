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

	zk *server.Server
}

func NewClient(username string, ip string, port uint) *Client {
	return &Client{
		User: chat_message.User{
			Username: username,
			Addr: &chat_message.Address{
				Ip:   ip,
				Port: uint32(port),
			},
		},
		peers: make(map[string]Peer),
	}
}

func (c *Client) ConnectZk(zk *server.Server) {
	c.zk = zk
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
		in := strings.Split(input, " ")
		if len(in) < 2 {
			fmt.Println("input: <command> <target>")
			continue
		}

		command := in[0]
		target := in[1]

		if command == "send" {
			go c.BroadcastMessage(target)
		} else if command == "create" {
			c.CreateChannel(target)
		} else if command == "getchans" {
			fmt.Println(c.ListChannels())
		} else if command == "join" {
			c.JoinChannel(target)
		} else if command == "conn" {
			users, _ := c.GetChannelUsers(target)
			for _, user := range users {
				if user.Username != c.User.Username {
					go c.DialAddress(user.Addr)
				}
			}
		} else if command == "sign" && target == "out" {
			usernames := make([]string, 0)
			for _, peer := range c.peers {
				usernames = append(usernames, peer.user.Username)
			}

			for _, username := range usernames {
				if err := c.CloseConnection(username); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
