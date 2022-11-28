package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	"strconv"
	"log"
	"fmt"
	"os"
	"bufio"
	"path/filepath"
	"time"
)

const (
	usersPath = "/users"
	connPath = "/conn"
	configFilePath = "../server/id.txt"
)

type Server struct {
	conn *zk.Conn
}

func (s *Server) Init(serverAddress string) {
	conn, _, err := zk.Connect([]string{serverAddress}, time.Second)
	s.conn = conn
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) RegisterUser(name string, ipv4 string, publicKey string) (string, error) {
	usersExists := zookeeper.CheckZNode(s.conn, usersPath)
	if usersExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", usersPath)
	}

	numberOfUsersString, version := zookeeper.GetZNode(s.conn, usersPath)
	numberOfUsersUpdated, _ := strconv.Atoi(numberOfUsersString)
	numberOfUsersUpdated++
	zookeeper.SetZNode(s.conn, usersPath, strconv.Itoa(numberOfUsersUpdated), version)

	userPath := fmt.Sprintf("%s/id%d", usersPath, numberOfUsersUpdated)
	userData := fmt.Sprintf("name %s\nipv4 %s\npublic-key %s", name, ipv4, publicKey)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(s.conn, userPath, flagPermanent, userData)
	CreateIdLocal(numberOfUsersUpdated)
	return ZNodePath, err
}

func (s *Server) SetUserOnline(userNumber int) (string, error) {
	connExists := zookeeper.CheckZNode(s.conn, connPath)
	if connExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", connPath)
	}

	userConnPath := fmt.Sprintf("%s/id%d", connPath, userNumber)
	ZNodeConnPath, err := zookeeper.CreateZNode(s.conn, userConnPath, zk.FlagEphemeral, "")
	return ZNodeConnPath, err
}

func GetIdFromLocal() (int, error) {
	id := 0
	absPath, absError := filepath.Abs(configFilePath)
	if absError != nil {
		log.Fatal(absError)
	}
	file, err := os.OpenFile(absPath, os.O_RDONLY, 0664)
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
	absPath, absError := filepath.Abs(configFilePath)
	if absError != nil {
		log.Fatal(absError)
	}
	file, openError := os.OpenFile(absPath, os.O_RDWR | os.O_CREATE, 0664)
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