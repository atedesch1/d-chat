package chat

import (
	"context"
	"fmt"
	"time"

	chat_message "github.com/decentralized-chat/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) SendMessage(ctx context.Context, msg *chat_message.ContentMessage) (*chat_message.AckMessage, error) {
	fmt.Println(msg.From.Username, ":", msg.Content)

	return &chat_message.AckMessage{
		From:   &c.User,
		SentAt: timestamppb.Now(),
	}, nil
}

func (c *Client) MessagePeer(msg *chat_message.ContentMessage, peer Peer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := peer.client.SendMessage(ctx, msg)
	if err != nil {
		return err
	}
	// log.Println("Ack from", ack.From.Username)

	return nil
}

func (c *Client) BroadcastMessage(content string) error {
	msg := &chat_message.ContentMessage{
		From:    &c.User,
		Content: content,
		SentAt:  timestamppb.Now(),
	}

	for _, peer := range c.peers {
		go c.MessagePeer(msg, peer)
	}

	return nil
}
