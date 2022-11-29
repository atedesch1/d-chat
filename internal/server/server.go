package server

import (
	"github.com/decentralized-chat/pkg/zookeeper"
	"github.com/go-zookeeper/zk"
	"strconv"
	"log"
	"fmt"
	"time"
	"strings"
	"errors"
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

func (ci *ChannelInfo) Init(channelname string, users []string) {
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
		return errors.New("Error when connecting to ZooKeeper.")
	}
	return nil
}

func ParseUserData(data string) *UserInfo {
	lines := strings.Split(data, "\n")
	username := strings.Split(lines[0], " ")[1]
	ipv4 := strings.Split(lines[1], " ")[1]
	port := strings.Split(lines[2], " ")[1]
	publicKey := strings.Split(lines[3], " ")[1]
	ui := new(UserInfo)
	ui.Init(username, ipv4, port, publicKey)
	return ui
}

func ParseChannelData(data string) *ChannelInfo {
	temp := strings.Split(data, "\n")
	channelname := strings.Split(temp[0], " ")[1]
	users := strings.Split(temp[1], " ")[1:]
	ci := new(ChannelInfo)
	ci.Init(channelname, users)
	return ci
}

func (s *Server) GetUserIdFromUsername(user string) (int, error) {
	children, _, err := s.conn.Children(usersPath)
	if err != nil {
		log.Fatal(err)
		return -1, errors.New("error when accessing /users children")
	}
	for _, userId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", usersPath, userId))
		ui := ParseUserData(data)
		if ui.username == user {
			userIdConverted, _ := strconv.Atoi(userId[2:])
			return userIdConverted, nil
		}
	}
	return -1, errors.New("username not found")
}

func (s *Server) RegisterUser(user string, ipv4 string, port string, publicKey string) error {
	usersExists := zookeeper.CheckZNode(s.conn, usersPath)
	if usersExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", usersPath)
		return errors.New("no path /users in zookeeper")
	}

	numberOfUsersString, version := zookeeper.GetZNode(s.conn, usersPath)
	numberOfUsersUpdated, _ := strconv.Atoi(numberOfUsersString)
	numberOfUsersUpdated++
	zookeeper.SetZNode(s.conn, usersPath, strconv.Itoa(numberOfUsersUpdated), version)

	userPath := fmt.Sprintf("%s/id%d", usersPath, numberOfUsersUpdated)
	userData := fmt.Sprintf("username %s\nipv4 %s\nport %s\npublic-key %s", user, ipv4, publicKey)
	flagPermanent := int32(0)
	_, err := zookeeper.CreateZNode(s.conn, userPath, flagPermanent, userData)
	return err
}

func (s *Server) SetUserOnline(user string) error {
	userId, err := s.GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	connExists := zookeeper.CheckZNode(s.conn, connPath)
	if connExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", connPath)
		return errors.New("no path /conn in zookeeper")
	}

	userConnPath := fmt.Sprintf("%s/id%d", connPath, userId)
	_, err = zookeeper.CreateZNode(s.conn, userConnPath, zk.FlagEphemeral, "")
	return err
}

func (s *Server) IsUserRegistered(user string) bool {
	userId, err := s.GetUserIdFromUsername(user)
	if err != nil {
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

func (s *Server) IsUserOnline(user string) (bool, error) {
	userExists := s.IsUserRegistered(user)
	if userExists == false {
		return false, errors.New("user was not registered")
	}
	children, _, err := s.conn.Children(connPath)
	if err != nil {
		log.Fatal(err)
	}
	userId, err := s.GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	for _, id := range children {
		userIdInt, _ := strconv.Atoi(id[2:])
		if userId == userIdInt {
			return true, nil
		}
	}
	return false, nil
}

func (s *Server) RegisterChannel(channelName string, user string) error {
	userId, err := s.GetUserIdFromUsername(user)
	if err != nil {
		log.Fatal(err)
	}
	channelsExists := zookeeper.CheckZNode(s.conn, channelsPath)
	if channelsExists == false {
		log.Fatalf("You must set %s path in the ZooKeeper.", channelsPath)
		return errors.New("no path /channels in zookeeper")
	}

	numberOfChannelsString, version := zookeeper.GetZNode(s.conn, channelsPath)
	numberOfChannelsUpdated, _ := strconv.Atoi(numberOfChannelsString)
	numberOfChannelsUpdated++
	zookeeper.SetZNode(s.conn, channelsPath, strconv.Itoa(numberOfChannelsUpdated), version)

	channelPath := fmt.Sprintf("%s/ch%d", channelsPath, numberOfChannelsUpdated)
	channelData := fmt.Sprintf("channelname %s\nusers id%d", channelName, userId)
	flagPermanent := int32(0)
	_, err = zookeeper.CreateZNode(s.conn, channelPath, flagPermanent, channelData)
	return err
}

func (s *Server) AddUserToChannel(channelName string, user string) error {
	userInsideChannel := s.IsUserInsideChannel(channelName, user)
	if userInsideChannel == true {
		return errors.New("user is already in the channel")
	}
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		ci := ParseChannelData(data)
		if ci.channelname == channelName {
			ci.users = append(ci.users, fmt.Sprintf(" %s\n", user))
			channelDataStr := GenerateChannelData(ci.channelname, ci.users)
			zookeeper.SetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId), channelDataStr, version)
			return nil
		}
	}
	return errors.New("channel does not exist")
}

func (s *Server) GetChannelUsers(channelname string) []string {
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		ci := ParseChannelData(data)
		if ci.channelname == channelname {
			return ci.users
		}
	}
	return []string{}
}

func (s *Server) GetChannelsName() []string {
	var channels []string
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, _ := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
		ci := ParseChannelData(data)
		channels = append(channels, ci.channelname)		
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
		ci := ParseChannelData(data)
		if ci.channelname == channelname {
			deletePath := fmt.Sprintf("%s/%s", channelsPath, channelId)
			zookeeper.DeleteZNode(s.conn, deletePath, version)
			return nil
		}
	}
	return errors.New("cannot delete a channel that does not exist")
}

func (s *Server) DeleteUserFromChannel(channelname string, user string) error { 
	children, _, err := s.conn.Children(channelsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, channelId := range children {
		data, version := zookeeper.GetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId))
	 	ci := ParseChannelData(data)
	 	if ci.channelname == channelname {
			newData := fmt.Sprintf("channelname %s\nusers", ci.channelname)
			for _, currUser := range ci.users {
				if currUser == user {
					continue
				}
				newData += fmt.Sprintf(" %s", currUser)
			}
			zookeeper.SetZNode(s.conn, fmt.Sprintf("%s/%s", channelsPath, channelId), newData, version)
			return nil
		}
	}
	return errors.New("cannot delete a user that is not in the channel")
}

func GenerateChannelData(channelname string, users []string) string {
	data := fmt.Sprintf("channelname %s\nusers", channelname)
	for _, username := range users {
		formatUsername := fmt.Sprintf(" %s", username)
		data += formatUsername
	}
	return data
}