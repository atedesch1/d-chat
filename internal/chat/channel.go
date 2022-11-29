package chat

import (
	"errors"
	"fmt"
	"strconv"

	chat_message "github.com/decentralized-chat/pb"
)

func (c *Client) RegisterUser() {
	if !c.zk.IsUserRegistered(c.User.Username) {
		c.zk.RegisterUser(c.User.Username, c.User.Addr.Ip, strconv.Itoa(int(c.User.Addr.Port)), "")
	}
	c.zk.SetUserOnline(c.User.Username)
}

func (c *Client) ListChannels() []string {
	return c.zk.GetChannelsName()
}

func (c *Client) CreateChannel(name string) error {
	channels := c.ListChannels()

	for _, channel := range channels {
		if channel == name {
			return errors.New("channel already exists")
		}
	}

	return c.zk.RegisterChannel(name, c.User.Username)
}

func (c *Client) RetrieveLostMessages() {
	messages, _ := c.zk.GetMessageFromQueue(c.User.Username)

	for _, message := range messages {
		fmt.Printf("%s: %s\n", message.From, message.Content)
	}
}

func (c *Client) JoinChannel(name string) {
	c.zk.AddUserToChannel(name, c.User.Username)
	users, _ := c.GetChannelUsers(name)
	for _, user := range users {
		if user.Username != c.User.Username {
			go c.DialAddress(user.Addr)
		}
	}
	c.RetrieveLostMessages()
	c.channel = name
}

func (c *Client) GetChannelUsers(name string) ([]*chat_message.User, error) {
	usernames := c.zk.GetChannelUsers(name)

	users := make([]*chat_message.User, len(usernames))
	for idx, username := range usernames {
		userInfo, err := c.zk.GetUserData(username)
		if err != nil {
			return users, err
		}

		port, _ := strconv.Atoi(userInfo.Port)

		addr := &chat_message.Address{
			Ip:   userInfo.Ipv4,
			Port: uint32(port),
		}

		users[idx] = &chat_message.User{
			Username: username,
			Addr:     addr,
		}
	}

	return users, nil
}

func (c *Client) DisconnectFromChannel() {
	usernames := make([]string, 0)
	for _, peer := range c.peers {
		usernames = append(usernames, peer.user.Username)
	}

	for _, username := range usernames {
		if err := c.CloseConnection(username); err != nil {
			fmt.Println(err)
		}
	}

	c.channel = ""
}

func (c *Client) LeaveChannel(name string) {
	c.zk.DeleteUserFromChannel(name, c.User.Username)
}

func (c *Client) DeleteChannel(name string) {
	c.zk.DeleteChannel(name)
}
