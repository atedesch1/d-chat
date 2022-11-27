package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"fmt"
	"strconv"
)

var _ = Describe("Server", func() {
	var clientId int

	When("A new user open the client", func() {
		It(`Should be able to register its 
		    information in the raw ZooKeeper server`, func() {
			local := "127.0.0.1"
			conn, _, err := zk.Connect([]string{local}, time.Second)
			Expect(err).To(BeNil())

			username := "username"
			ipv4 := "192.168.0.10"
			publicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
			path, err := RegisterUser(conn, username, ipv4, publicKey)
			usersData, _ := zookeeper.GetZNode(conn, usersPath)
			clientId, _ = strconv.Atoi(usersData)
			expectedPath := fmt.Sprintf("%s/id%s", usersPath, usersData)
			Expect(err).To(BeNil())
			Expect(path).To(Equal(expectedPath))
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
})