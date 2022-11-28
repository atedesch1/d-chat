package chat

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

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
		peers: make(map[string]Peer),
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
		in := strings.Split(input, " ")
		if len(in) < 2 {
			fmt.Println("input: <command> <target>")
			continue
		}

		command := in[0]
		target := in[1]

		if command == "send" {
			go c.BroadcastMessage(target)
		} else if command == "conn" {
			port, err := strconv.Atoi(target)
			if err != nil {
				log.Println("couldnt convert porn")
				continue
			}

			go c.DialAddress(&chat_message.Address{
				Ip:   "localhost",
				Port: uint32(port),
			})
		} else if command == "disc" {
			if err := c.CloseConnection(target); err != nil {
				fmt.Println(err)
			}
		}
	}
}
