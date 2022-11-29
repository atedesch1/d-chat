package server

import (
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

			channelname := "mychannel"
			err = s.RegisterChannel(channelname, username)
			Expect(err).To(BeNil())
			status = s.IsUserInsideChannel(channelname, username)
			Expect(status).To(Equal(true))

			err = s.DeleteChannel(channelname)
			Expect(err).To(BeNil())
		})
	})

	// var clientId int

	// When("A new user open the client", func() {
	// 	server := new(Server)
	// 	server.Init("127.0.0.1:2181")

	// 	It(`Should be able to register its 
	// 	    information in the raw ZooKeeper server`, func() {
	// 		username := "username"
	// 		ipv4 := "192.168.0.10"
	// 		publicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
	// 		path, err := server.RegisterUser(username, ipv4, publicKey)
	// 		usersData, _ := zookeeper.GetZNode(server.conn, usersPath)
	// 		clientId, _ = strconv.Atoi(usersData)
	// 		expectedPath := fmt.Sprintf("%s/id%s", usersPath, usersData)
	// 		currentConnPath, connErr := server.SetUserOnline(clientId)
	// 		expectedConnPath := fmt.Sprintf("%s/id%d", connPath, clientId)
	// 		Expect(err).To(BeNil())
	// 		Expect(connErr).To(BeNil())
	// 		Expect(path).To(Equal(expectedPath))
	// 		Expect(currentConnPath).To(Equal(expectedConnPath))
	// 	})

	// 	It(`Should be able to get the user's id from the user's name`, func() {
	// 		userId := server.GetUserIdFromUsername("username")
	// 		Expect(userId).NotTo(Equal(-1))
	// 	})
	// })

	// When("A user creates a new channel", func() {
	// 	server := new(Server)
	// 	server.Init("127.0.0.1:2181")
	// 	userId := 1

	// 	It(`Should stores the channel name and
	// 		the users connected to it`, func() {
	// 		channelName := "channelname"
	// 		channelPath, err := server.RegisterChannel(channelName, userId)
	// 		channelData, _ := zookeeper.GetZNode(server.conn, channelsPath)
	// 		expectedChannelPath := fmt.Sprintf("%s/ch%s", channelsPath, channelData)
	// 		Expect(err).To(BeNil())
	// 		Expect(channelPath).To(Equal(expectedChannelPath))
	// 	})

	// 	It(`Should be able to delete this channel`, func() {
	// 		channelName := "deletedchannel"
	// 		server.RegisterChannel(channelName, userId)
	// 		status := server.DeleteChannel(channelName)
	// 		Expect(status).To(Equal(true))
	// 	})

	// 	It(`Should be able to insert new users in the channel`, func() {
	// 		channelName := "insertchannel"
	// 		server.RegisterChannel(channelName, userId)
	// 		var newUsersId []int
	// 		newUsersId = append(newUsersId, 90)
	// 		newUsersId = append(newUsersId, 40)
	// 		newUsersId = append(newUsersId, 50)
	// 		newUsersId = append(newUsersId, 60)
	// 		statusAdd := server.AddUsersToChannel(channelName, newUsersId)
	// 		expectedUsersId := [5]int{userId, 90, 40, 50, 60}
	// 		channelUsersId := server.GetChannelUsers(channelName)
	// 		fmt.Print(channelUsersId)
	// 		fmt.Print(expectedUsersId)
	// 		for index, _ := range channelUsersId {
	// 			Expect(channelUsersId[index]).To(Equal(expectedUsersId[index]))
	// 		}
	// 		Expect(statusAdd).To(Equal(true))
	// 		statusDelete := server.DeleteChannel(channelName)
	// 		Expect(statusDelete).To(Equal(true))
	// 	})

	// 	It("Should be able to be removed", func() {
	// 		channelName := "insert-remove-channel"
	// 		usernames := [3]string{"mike", "paul", "john"}
	// 		ipv4 := [3]string{"ip1", "ip2", "ip3"}
	// 		publicKey := [3]string{"key1", "key2", "key3"}
	// 		var clientIds []int
	// 		for index, _ := range usernames {
	// 			server.RegisterUser(usernames[index], ipv4[index], publicKey[index])
	// 			usersData, _ := zookeeper.GetZNode(server.conn, usersPath)
	// 			clientId, _ = strconv.Atoi(usersData)
	// 			clientIds = append(clientIds, clientId)
	// 		}
	// 		server.RegisterChannel(channelName, clientIds[0])
	// 		statusAdd := server.AddUsersToChannel(channelName, clientIds[1:])
	// 		channelUsers := server.GetChannelUsers(channelName)
	// 		fmt.Println(channelUsers)
	// 		server.DeleteUserFromChannel(channelName, "paul")
	// 		newChannelUsers := server.GetChannelUsers(channelName)
	// 		fmt.Println(newChannelUsers)
	// 		Expect(statusAdd).To(Equal(true))
	// 		statusDelete := server.DeleteChannel(channelName)
	// 		Expect(statusDelete).To(Equal(true))
	// 	})
	// })
	
	// It(`Should parse the data correctly`, func() {
	// 	data := "channel-name channel\nusers id0 id1 id2 id100 id90 id45"
	// 	channelName, idList := ParseChannelData(data)
	// 	channelNameExpected := "channel"
	// 	Expect(channelName).To(Equal(channelNameExpected))
	// 	idListExpected := [6]int{0, 1, 2, 100, 90, 45}
	// 	for index, _ := range idListExpected {
	// 		Expect(idList[index]).To(Equal(idListExpected[index]))
	// 	}
	// })
	
	// It(`Should parse the data correctly`, func() {
	// 	expectedUsername := "john.doe"
	// 	expectedIpv4 := "192.168.0.10"
	// 	expectedPublicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
	// 	data := fmt.Sprintf("name %s\nipv4 %s\npublic-key %s", expectedUsername, expectedIpv4, expectedPublicKey)
	// 	username, ipv4, publicKey := ParseUserData(data)
	// 	Expect(username).To(Equal(expectedUsername))
	// 	Expect(ipv4).To(Equal(expectedIpv4))
	// 	Expect(publicKey).To(Equal(expectedPublicKey))
	// })
})