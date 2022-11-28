package server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/decentralized-chat/pkg/zookeeper"
	"fmt"
	"strconv"
)

var _ = Describe("Server", func() {
	var clientId int

	When("A new user open the client", func() {
		It(`Should be able to register its 
		    information in the raw ZooKeeper server`, func() {
			server := new(Server)
			server.Init("127.0.0.1:2181")

			username := "username"
			ipv4 := "192.168.0.10"
			publicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
			path, err := server.RegisterUser(username, ipv4, publicKey)
			usersData, _ := zookeeper.GetZNode(server.conn, usersPath)
			clientId, _ = strconv.Atoi(usersData)
			expectedPath := fmt.Sprintf("%s/id%s", usersPath, usersData)
			currentConnPath, connErr := server.SetUserOnline(clientId)
			expectedConnPath := fmt.Sprintf("%s/id%d", connPath, clientId)
			Expect(err).To(BeNil())
			Expect(connErr).To(BeNil())
			Expect(path).To(Equal(expectedPath))
			Expect(currentConnPath).To(Equal(expectedConnPath))
		})
	})

	When("A old user connects", func() {
		It(`Should be able to get its id
		    from the local storage`, func() {
			id, idError := GetIdFromLocal()
			Expect(idError).To(BeNil())
			Expect(id).To(Equal(clientId))
		})
	})

	When("A user creates a new channel", func() {
		It(`Should stores the channel name and
			the users connected to it`, func() {
			server := new(Server)
			server.Init("127.0.0.1:2181")

			channelName := "channelname"
			userId := 1
			channelPath, err := server.RegisterChannel(channelName, userId)
			channelData, _ := zookeeper.GetZNode(server.conn, channelsPath)
			expectedChannelPath := fmt.Sprintf("%s/ch%s", channelsPath, channelData)
			Expect(err).To(BeNil())
			Expect(channelPath).To(Equal(expectedChannelPath))
		})
	})
})