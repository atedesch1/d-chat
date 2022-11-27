package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"log"
)

func checkZNode(conn *zk.Conn, zkPath string) bool {
	exists, _, err := conn.Exists(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func getZNode(conn *zk.Conn, zkPath string) string {
	data, _, err := conn.Get(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data[:])
	return s
}

func registerUser(conn *zk.Conn, zkPath string, zkFlags int32, data string) (string, error) {
	exists := checkZNode(conn, zkPath)
	if exists == false {
		create, err := conn.Create(zkPath, []byte(data), zkFlags, zk.WorldACL(zk.PermAll))
		return create, err
	}
	return zkPath, nil
}