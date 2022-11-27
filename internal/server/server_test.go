package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"fmt"
)

var _ = Describe("Server", func() {
	When("A new user open the client", func() {
		It(`Should be able to register its 
		    information in the raw ZooKeeper server`, func() {
			local := "127.0.0.1"
			conn, _, err := zk.Connect([]string{local}, time.Second)
			Expect(err).To(BeNil())

			username := "username"
			ipv4 := "192.168.0.10"
			publicKey := "i1RfARNCYn9+K3xmRNTaXG9sVSK6TMgY9l8SDm3MUZ4="
			path, err := zookeeper.RegisterUser(conn, username, ipv4, publicKey)
			usersData, _ := zookeeper.GetZNode(conn, zookeeper.usersPath)
			expectedPath := fmt.Sprintf("%s/id%s", zookeeper.usersPath, usersData)
			Expect(err).To(BeNil())
			Expect(path).To(Equal(expectedPath))
		})
	})
})