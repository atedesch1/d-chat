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
	"strings"
)

const (
	usersPath = "/users"
	connPath = "/conn"
	channelsPath = "/channels"
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

func (s *Server) SetUserOnline(userId int) (string, error) {
	connExists := zookeeper.CheckZNode(s.conn, connPath)
	if connExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", connPath)
	}

	userConnPath := fmt.Sprintf("%s/id%d", connPath, userId)
	ZNodeConnPath, err := zookeeper.CreateZNode(s.conn, userConnPath, zk.FlagEphemeral, "")
	return ZNodeConnPath, err
}

func (s *Server) RegisterChannel(channelName string, userId int) (string, error) {
	channelsExists := zookeeper.CheckZNode(s.conn, channelsPath)
	if channelsExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", channelsPath)
	}

	numberOfChannelsString, version := zookeeper.GetZNode(s.conn, channelsPath)
	numberOfChannelsUpdated, _ := strconv.Atoi(numberOfChannelsString)
	numberOfChannelsUpdated++
	zookeeper.SetZNode(s.conn, channelsPath, strconv.Itoa(numberOfChannelsUpdated), version)

	channelPath := fmt.Sprintf("%s/ch%d", channelsPath, numberOfChannelsUpdated)
	channelData := fmt.Sprintf("channel-name %s\nusers id%d", channelName, userId)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(s.conn, channelPath, flagPermanent, channelData)
	return ZNodePath, err
}

func (s *Server) AddUsersToChannel(channelName string, idList []int) bool {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, channelData := ParseChannelData(data)
		if currChannelName == channelName {
			for _, id := range idList {
				channelData = append(channelData, id)
			}
			channelDataStr := GenerateChannelData(currChannelName, channelData)
			zookeeper.SetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId), channelDataStr, version)
			return true
		}
	}
	return false
}

func (s *Server) GetChannelUsers(channelName string) []int {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, channelData := ParseChannelData(data)
		if currChannelName == channelName {
			return channelData
		}
	}
	return nil
}

func (s *Server) GetChannelsName() []string {
	var channels []string
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, _ := ParseChannelData(data)
		channels = append(channels, currChannelName)		
	}
	return channels
}

func (s *Server) DeleteChannel(channelName string) bool {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, _ := ParseChannelData(data)
		if currChannelName == channelName {
			deletePath := fmt.Sprintf("%s/%s", channelsPath, channelId)
			zookeeper.DeleteZNode(s.conn, deletePath, version)
			return true
		}
	}
	return false
}

func (s *Server) DeleteUserFromChannel(channelName string, user string) { }

func GenerateChannelData(channelName string, idList []int) string {
	data := fmt.Sprintf("channel-name %s\nusers", channelName)
	for _, id := range idList {
		newId := fmt.Sprintf(" id%d", id)
		data += newId
	}
	return data
}

func ParseUserData(data string) (string, string, string) {
	lines := strings.Split(data, "\n")
	username := strings.Split(lines[0], " ")[1]
	ipv4 := strings.Split(lines[1], " ")[1]
	publicKey := strings.Split(lines[2], " ")[1]
	return username, ipv4, publicKey
}

func ParseChannelData(data string) (string, []int) {
	temp := strings.Split(data, "\n")
	channelName := strings.Split(temp[0], " ")[1]
	
	var idList []int
	idStringList := strings.Split(temp[1], " ")[1:]
	for _, idStr := range idStringList {
		id, err := strconv.Atoi(idStr[2:])
		if err != nil {
			log.Fatal(err)
		}
		idList = append(idList, id)
	}
	return channelName, idList
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