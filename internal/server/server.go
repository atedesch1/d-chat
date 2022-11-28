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
)

type Server struct {
	conn *zk.Conn
}

type UserInfo struct {
	username  string
	ipv4 	  string
	port 	  string
	publicKey string
}

type ChannelInfo struct {
	channelname   string
	users       []string
}

func (ci *ChannelInfo) Init(channelname string, users []string) struct {
	ci.channelname = channelname
	ci.users = users
}

func (ui *UserInfo) Init(username string, ipv4 string, port string, publicKey string) {
	ui.username = username
	ui.ipv4 = ipv4
	ui.port = port
	ui.publicKey = publicKey
}

func (s *Server) Init(ipv4 string, port string) error {
	addr := fmt.Sprintf("%s:%s", ipv4, port)
	conn, _, err := zk.Connect([]string{addr}, time.Second)
	s.conn = conn
	if err != nil {
		log.Fatal(err)
		return error.New("Error when connecting to ZooKeeper.")
	}
	return nil
}

func (s *Server) RegisterUser(user string, ipv4 string, port string, publicKey string) error {
	usersExists := zookeeper.CheckZNode(s.conn, usersPath)
	if usersExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", usersPath)
		return error.New("no path /users in zookeeper")
	}

	numberOfUsersString, version := zookeeper.GetZNode(s.conn, usersPath)
	numberOfUsersUpdated, _ := strconv.Atoi(numberOfUsersString)
	numberOfUsersUpdated++
	zookeeper.SetZNode(s.conn, usersPath, strconv.Itoa(numberOfUsersUpdated), version)

	userPath := fmt.Sprintf("%s/id%d", usersPath, numberOfUsersUpdated)
	userData := fmt.Sprintf("username %s\nipv4 %s\nport %s\npublic-key %s", user, ipv4, publicKey)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(s.conn, userPath, flagPermanent, userData)
	return err
}

func (s *Server) SetUserOnline(user string) (string, error) error {
	userId, err := GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	connExists := zookeeper.CheckZNode(s.conn, connPath)
	if connExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", connPath)
		return error.New("no path /conn in zookeeper")
	}

	userConnPath := fmt.Sprintf("%s/id%d", connPath, userId)
	ZNodeConnPath, err := zookeeper.CreateZNode(s.conn, userConnPath, zk.FlagEphemeral, "")
	return err
}

func (s *Server) IsUserRegistered(user string) bool {
	userId, err := GetUserIdFromUsername(user)
	err != nil {
		log.Fatal(err)
	}
	userExists := zookeeper.CheckZNode(s.conn, fmt.Sprintf("%s/id%d", usersPath, userId))
	return userExists
}

func (s *Server) IsUserInsideChannel(channelname string, user string) bool {
	users := s.GetChannelUsers(channelname)
	for _, currUser := range users {
		if user == currUser {
			return true
		}
	}
	return false
}

func (s *Server) IsUserOnline(user string) bool, error {
	userExists := IsUserRegistered(user)
	if userExists == false {
		return false, error.New("user was not registered")
	}
	children, _, err := s.conn.Children(connPath)
	if err != nil {
		log.Fatal(err)
	}
	userId, err := GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	for _, id := range children {
		if userId == strconv.Atoi(id[2:]) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Server) RegisterChannel(channelName string, user string) error {
	userId, err := GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	channelsExists := zookeeper.CheckZNode(s.conn, channelsPath)
	if channelsExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", channelsPath)
		return error.New("no path /channels in zookeeper")
	}

	numberOfChannelsString, version := zookeeper.GetZNode(s.conn, channelsPath)
	numberOfChannelsUpdated, _ := strconv.Atoi(numberOfChannelsString)
	numberOfChannelsUpdated++
	zookeeper.SetZNode(s.conn, channelsPath, strconv.Itoa(numberOfChannelsUpdated), version)

	channelPath := fmt.Sprintf("%s/ch%d", channelsPath, numberOfChannelsUpdated)
	channelData := fmt.Sprintf("channelname %s\nusers id%d", channelName, userId)
	flagPermanent := int32(0)
	ZNodePath, err := zookeeper.CreateZNode(s.conn, channelPath, flagPermanent, channelData)
	return err
}

func (s *Server) AddUserToChannel(channelName string, user string) error {
	userInsideChannel := IsUserInsideChannel(channelname, user)
	if userInsideChannel == true {
		return error.New("user is already in the channel")
	}
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, channelData := ParseChannelData(data)
		if currChannelName == channelName {
			channelData = append(channelData, fmt.Sprintf(" %s\n", user))
			channelDataStr := GenerateChannelData(currChannelName, channelData)
			zookeeper.SetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId), channelDataStr, version)
			return nil
		}
	}
	return error.New("channel does not exist")
}

func (s *Server) GetChannelUsers(channelname string) []string {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, channelData := ParseChannelData(data)
		if currChannelName == channelname {
			return channelData
		}
	}
	return []string
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

func (s *Server) DeleteChannel(channelname string) error {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		currChannelName, _ := ParseChannelData(data)
		if currChannelName == channelname {
			deletePath := fmt.Sprintf("%s/%s", channelsPath, channelId)
			zookeeper.DeleteZNode(s.conn, deletePath, version)
			return nil
		}
	}
	return error.New("cannot delete a channel that does not exist")
}

func (s *Server) DeleteUserFromChannel(channelname string, user string) error { 
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
	 	currChannelName, currUsers := ParseChannelData(data)
	 	if currChannelName == channelname {
			userId := s.GetUserIdFromUsername(user)
			newData := fmt.Sprintf("channelname %s\nusers", currChannelName)
			for _, currUser := range currUsers {
				if currUser == user {
					continue
				}
				newData += fmt.Sprintf(" %s", currUser)
			}
			zookeeper.SetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId), newData, version)
			return nil
		}
	}
	return error.New("cannot delete a user that is not in the channel")
}

func (s *Server) GetUserIdFromUsername(user string) (int, error) {
	children, _, err := s.conn.Children(usersPath)
	if err != nil {
		log.Fatal(err)
		return -1, error.New("error when accessing /users children")
	}
	for _, userId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", usersPath, userId))
		username, _, _ := ParseUserData(data)
		if username == user {
			userIdConverted, _ := strconv.Atoi(userId[2:])
			return userIdConverted, nil
		}
	}
	return -1, error.New("username not found")
}

func GenerateChannelData(channelname string, users []string) string {
	data := fmt.Sprintf("channelname %s\nusers", channelname)
	for _, username := range users {
		formatUsername := fmt.Sprintf(" %s", username)
		data += formatUsername
	}
	return data
}

func ParseUserData(data string) UserInfo {
	lines := strings.Split(data, "\n")
	username := strings.Split(lines[0], " ")[1]
	ipv4 := strings.Split(lines[1], " ")[1]
	port := strings.Split(lines[2], " ")[1]
	publicKey := strings.Split(lines[3], " ")[1]
	ui := New(UserInfo)
	ui.Init(username, ipv4, port, publicKey)
	return ui
}

func ParseChannelData(data string) ChannelInfo {
	temp := strings.Split(data, "\n")
	channelname := strings.Split(temp[0], " ")[1]
	users := strings.Split(tempo[1], " ")[1:]
	ci := new(ChannelInfo)
	ci.Init(channelname, users)
	return ci
}