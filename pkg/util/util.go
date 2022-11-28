package util

import (
	"fmt"
)

func JoinIpAndPort(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
