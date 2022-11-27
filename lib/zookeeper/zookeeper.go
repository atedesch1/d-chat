package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"fmt"
	"time"
	"log"
)

func checkZNode(zkConn *zk.Conn, zkPath string) bool {
	exists, _, err := zkConn.Exists(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func getZNode(zkConn *zk.Conn, zkPath string) string {
	data, _, err := zkConn.Get(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data[:])
	return s
}

func registerUser(zkConn *zk.Conn, zkPath string, zkFlags int, data string) {
	exists := checkZNode(zkConn, zkPath)
	if exists == false {
		zkConn.Create(zkPath, []byte(data), zkFlags, zk.WorldACL(zk.PermAll))
	}
}