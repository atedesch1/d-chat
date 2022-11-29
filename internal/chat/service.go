package chat

import (
	"context"
	"errors"
	"fmt"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) SendMessage(
	ctx context.Context,
	msg *chat_message.ContentMessage) (*chat_message.AckMessage, error) {
	fmt.Printf("%s: %s\n", msg.From.Username, msg.Content)
	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, nil
}

func (c *Client) GetUsername(
	ctx context.Context,
	in *emptypb.Empty) (*chat_message.UserInfo, error) {
	return &chat_message.UserInfo{
		User: &c.User,
	}, nil
}

func (c *Client) RequestConnection(
	ctx context.Context,
	msg *chat_message.ConnectionMessage) (*chat_message.AckMessage, error) {
	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, c.DialAddress(msg.User.Addr)
}

func (c *Client) Disconnect(
	ctx context.Context,
	msg *chat_message.ConnectionMessage) (*chat_message.AckMessage, error) {
	var err error

	if _, ok := c.peers[msg.User.Username]; ok {
		c.peers[msg.User.Username].conn.Close()
		c.RemovePeer(msg.User.Username)
	} else {
		err = errors.New("peer not connected")
	}

	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, err
}
