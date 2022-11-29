package util

import (
	"fmt"
	"strconv"
	"strings"

	chat_message "github.com/decentralized-chat/pb"
)

func JoinIpAndPort(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func HostToAddr(host string) *chat_message.Address {
	address := strings.Split(host, ":")
	ip := address[0]
	port, _ := strconv.Atoi(address[1])
	return &chat_message.Address{
		Ip:   ip,
		Port: uint32(port),
	}
}

func AddrToHost(addr *chat_message.Address) string {
	return JoinIpAndPort(addr.Ip, int(addr.Port))
}
