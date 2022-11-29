package server

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	s := new(Server)
	s.Init("127.0.0.1", "2181")

	username := "paul"
	ipv4 := "192.168.0.10"
	port := "2181"
	publicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
	channelname := "mychannel"

	user1 := "mike"
	user2 := "john"

	When("A channel data is retrieved from zk", func() {
		It("Should be able to parse it", func() {
			data := fmt.Sprintf("channelname %s\nusers %s %s %s", channelname, username, user1, user2)
			ci := ParseChannelData(data)
			expectedUsers := [...]string{username, user1, user2}
			for idx, _ := range expectedUsers {
				Expect(ci.users[idx]).To(Equal(expectedUsers[idx]))
			}
		})
	})

	When("We've got the channel name and the users", func() {
		It("Should be able to transform to zk data", func() {
			var users []string
			users = append(users, username, user1, user2)
			data := GenerateChannelData(channelname, users)
			expectedData := fmt.Sprintf("channelname %s\nusers %s %s %s", channelname, username, user1, user2)
			Expect(data).To(Equal(expectedData))
		})
	})

	When("A user data is retrieved from zk", func() {
		It("Should be able to parse it", func() {
			data := fmt.Sprintf("username %s\nipv4 %s\nport %s\npublic-key %s", username, ipv4, port, publicKey)
			ui := ParseUserData(data)
			Expect(ui.Username).To(Equal(username))
			Expect(ui.Ipv4).To(Equal(ipv4))
			Expect(ui.Port).To(Equal(port))
			Expect(ui.PublicKey).To(Equal(publicKey))
		})
	})

	When("A new user open the client", func() {
		It("Should be able to register its information", func() {
			status := s.IsUserRegistered(username)
			Expect(status).To(Equal(false))
			err := s.RegisterUser(username, ipv4, port, publicKey)
			Expect(err).To(BeNil())
			status = s.IsUserRegistered(username)
			Expect(status).To(Equal(true))
			err = s.SetUserOnline(username)
			Expect(err).To(BeNil())
			status, err = s.IsUserOnline(username)
			Expect(err).To(BeNil())
			Expect(status).To(Equal(true))
			ui, uiErr := s.GetUserData(username)
			Expect(uiErr).To(BeNil())
			Expect(ui.Username).To(Equal(username))
			Expect(ui.Ipv4).To(Equal(ipv4))
			Expect(ui.Port).To(Equal(port))
			Expect(ui.PublicKey).To(Equal(publicKey))
		})

		It("Should be able to create a channel", func() {
			err := s.RegisterChannel(channelname, username)
			Expect(err).To(BeNil())
			status := s.IsUserInsideChannel(channelname, username)
			Expect(status).To(Equal(true))
			chann := s.GetChannelsName()
			expectedChann := [...]string{channelname}
			for idx, _ := range chann {
				Expect(chann[idx]).To(Equal(expectedChann[idx]))
			}
		})

		It("Should be able to add users to a channel", func() {
			err := s.AddUserToChannel(channelname, user1)
			Expect(err).To(BeNil())
			status := s.IsUserInsideChannel(channelname, user1)
			Expect(status).To(Equal(true))
			err = s.AddUserToChannel(channelname, user2)
			Expect(err).To(BeNil())
			status = s.IsUserInsideChannel(channelname, user2)
			Expect(status).To(Equal(true))
			var expectedUsers1 []string
			expectedUsers1 = append(expectedUsers1, username, user1, user2)
			users1 := s.GetChannelUsers(channelname)
			Expect(users1).To(Equal(expectedUsers1))
			err = s.DeleteUserFromChannel(channelname, user1)
			Expect(err).To(BeNil())
			var expectedUsers2 []string
			expectedUsers2 = append(expectedUsers2, username, user2)
			users2 := s.GetChannelUsers(channelname)
			Expect(users2).To(Equal(expectedUsers2))
			err = s.SendMessageToQueue(username, user2, "message")
			Expect(err).To(BeNil())
			err = s.SendMessageToQueue(username, user2, "message message message")
			queue, queueErr := s.GetMessageFromQueue(username)
			Expect(queueErr).To(BeNil())
			Expect(queue[0].from).To(Equal(user2))
			Expect(queue[0].content).To(Equal("message"))
			Expect(queue[1].from).To(Equal(user2))
			Expect(queue[1].content).To(Equal("message message message"))
		})

		It("Should be able to delete a channel", func() {
			err := s.DeleteChannel(channelname)
			Expect(err).To(BeNil())
		})
	})
})
