package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	"strconv"
	"log"
	"fmt"
	"os"
	"bufio"
)

const (
	usersPath = "/users"
	connPath = "/conn"
	configFilePath = "/conf/id.txt"
)

func RegisterUser(conn *zk.Conn, name string, ipv4 string, publicKey string) (string, error) {
	usersExists := zookeeper.CheckZNode(conn, usersPath)
	if usersExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", usersPath)
	}

	numberOfUsersString, version := zookeeper.GetZNode(conn, usersPath)
	numberOfUsersUpdated, _ := strconv.Atoi(numberOfUsersString)
	numberOfUsersUpdated++
	zookeeper.SetZNode(conn, usersPath, strconv.Itoa(numberOfUsersUpdated), version)

	userPath := fmt.Sprintf("%s/id%d", usersPath, numberOfUsersUpdated)
	userData := fmt.Sprintf("name %s\nipv4 %s\npublic-key %s", name, ipv4, publicKey)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(conn, userPath, flagPermanent, userData)
	CreateIdLocal(numberOfUsersUpdated)
	return ZNodePath, err
}

func SetUserOnline(conn *zk.Conn, userNumber int) (string, error) {
	connExists := zookeeper.CheckZNode(conn, connPath)
	if connExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", connPath)
	}

	userConnPath := fmt.Sprintf("%s/id%d", connPath, userNumber)
	ZNodeConnPath, err := zookeeper.CreateZNode(conn, userConnPath, zk.FlagEphemeral, "")
	return ZNodeConnPath, err
}

func GetIdFromLocal() (int, error) {
	id := 0
	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0664)
	if err != nil {
		return id, err
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		id, _ = strconv.Atoi(fileScanner.Text())
	}
	return id, nil
}

func CreateIdLocal(id int) {
	file, openError := os.OpenFile(configFilePath, os.O_CREATE | os.O_RDWR, 0664)
	if openError != nil {
		log.Fatal(openError)
	}
	defer file.Close()
	idString := fmt.Sprintf("%d", id)
	_, writeError := file.WriteString(idString)
	if writeError != nil {
		log.Fatal(writeError)
	}
}