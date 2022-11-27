package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	"strconv"
	"log"
	"fmt"
)

const (
	usersPath = "/users"
)

func registerUser(conn *zk.Conn, name string, ipv4 string, publicKey string) (string, error) {
	usersExists := zookeeper.CheckZNode(conn, usersPath)
	if usersExists == false {
		log.Fatal("You must set /users path in the ZooKeeper.")
	}

	numberOfUsersString, version := zookeeper.GetZNode(conn, usersPath)
	numberOfUsersUpdated, _ := strconv.Atoi(numberOfUsersString)
	numberOfUsersUpdated++
	zookeeper.SetZNode(conn, usersPath, strconv.Itoa(numberOfUsersUpdated), version)

	userPath := fmt.Sprintf("%s/id%d", usersPath, numberOfUsersUpdated)
	userData := fmt.Sprintf("name\n%s\nipv4\n%s\npublic-key\n%s", name, ipv4, publicKey)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(conn, userPath, flagPermanent, userData)
	return ZNodePath, err
}
