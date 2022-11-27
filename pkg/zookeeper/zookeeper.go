package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"log"
)

func CheckZNode(conn *zk.Conn, zkPath string) bool {
	exists, _, err := conn.Exists(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func GetZNode(conn *zk.Conn, zkPath string) (string, int32) {
	data, stat, err := conn.Get(zkPath)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data[:])
	return s, stat.Version
}

func SetZNode(conn *zk.Conn, zkPath string, data string, version int32) *zk.Stat {
	stat, err := conn.Set(zkPath, []byte(data), version)
	if err != nil {
		log.Fatal(err)
	}
	return stat
}

func CreateZNode(conn *zk.Conn, zkPath string, zkFlags int32, data string) (string, error) {
	exists := CheckZNode(conn, zkPath)
	if exists == false {
		create, err := conn.Create(zkPath, []byte(data), zkFlags, zk.WorldACL(zk.PermAll))
		return create, err
	}
	return zkPath, nil
}