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
		server := new(Server)
		server.Init("127.0.0.1:2181")
		userId := 1

		It(`Should stores the channel name and
			the users connected to it`, func() {
			channelName := "channelname"
			channelPath, err := server.RegisterChannel(channelName, userId)
			channelData, _ := zookeeper.GetZNode(server.conn, channelsPath)
			expectedChannelPath := fmt.Sprintf("%s/ch%s", channelsPath, channelData)
			Expect(err).To(BeNil())
			Expect(channelPath).To(Equal(expectedChannelPath))
		})

		It(`Should be able to delete this channel`, func() {
			channelName := "deletedchannel"
			server.RegisterChannel(channelName, userId)
			status := server.DeleteChannel(channelName)
			Expect(status).To(Equal(true))
		})

		It(`Should be able to insert new users in the channel`, func() {
			channelName := "insertchannel"
			server.RegisterChannel(channelName, userId)
			var newUsersId []int
			newUsersId = append(newUsersId, 90)
			newUsersId = append(newUsersId, 40)
			newUsersId = append(newUsersId, 50)
			newUsersId = append(newUsersId, 60)
			statusAdd := server.AddUsersToChannel(channelName, newUsersId)
			expectedUsersId := [5]int{userId, 90, 40, 50, 60}
			channelUsersId := server.GetChannelUsers(channelName)
			fmt.Print(channelUsersId)
			fmt.Print(expectedUsersId)
			for index, _ := range channelUsersId {
				Expect(channelUsersId[index]).To(Equal(expectedUsersId[index]))
			}
			Expect(statusAdd).To(Equal(true))
			statusDelete := server.DeleteChannel(channelName)
			Expect(statusDelete).To(Equal(true))
		})
	})
	
	It(`Should parse the data correctly`, func() {
		data := "channel-name channel\nusers id0 id1 id2 id100 id90 id45"
		channelName, idList := ParseChannelData(data)
		channelNameExpected := "channel"
		Expect(channelName).To(Equal(channelNameExpected))
		idListExpected := [6]int{0, 1, 2, 100, 90, 45}
		for index, _ := range idListExpected {
			Expect(idList[index]).To(Equal(idListExpected[index]))
		}
	})
	
	It(`Should parse the data correctly`, func() {
		expectedUsername := "john.doe"
		expectedIpv4 := "192.168.0.10"
		expectedPublicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
		data := fmt.Sprintf("name %s\nipv4 %s\npublic-key %s", expectedUsername, expectedIpv4, expectedPublicKey)
		username, ipv4, publicKey := ParseUserData(data)
		Expect(username).To(Equal(expectedUsername))
		Expect(ipv4).To(Equal(expectedIpv4))
		Expect(publicKey).To(Equal(expectedPublicKey))
	})
})