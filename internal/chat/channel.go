package chat

import (
	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/grpc"
)

type Channel struct {
	id    uint
	users []*chat_message.User
}

type Peer struct {
	user   *chat_message.User
	client chat_message.ChatServiceClient
	conn   *grpc.ClientConn
}

func NewChannel(id uint) *Channel {
	return &Channel{
		id: id,
	}
}

func (c *Channel) AddUser(user *chat_message.User) {
	c.users = append(c.users, user)
}
