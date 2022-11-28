package chat

import (
	chat_message "github.com/decentralized-chat/pb"
)

type Channel struct {
	id    uint
	users []*chat_message.User
}

func NewChannel(id uint) *Channel {
	return &Channel{
		id: id,
	}
}

func (c *Channel) AddUser(user *chat_message.User) {
	c.users = append(c.users, user)
}
