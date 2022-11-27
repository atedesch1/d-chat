package main

import (
	"github.com/go-zookeeper/zk"
	"fmt"
	"time"
	"log"
)

func main() {
	// Constants
	// zk.FlagEphemeral
	// zk.FlagPermanent
	DefaultIp := "127.0.0.1"
	DefaultPort := "2181"
	zkPath := "/username"

	fmt.Println(DefaultIp)
	fmt.Println(DefaultPort)

	conn, _, err := zk.Connect([]string{DefaultIp}, time.Second)
	conn.Create(zkPath, []byte(DefaultIp), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Fatal(err)
	}
	children, stat, ch, err := conn.ChildrenW("/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)
}
