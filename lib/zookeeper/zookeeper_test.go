package zookeeper

import(
	"github.com/go-zookeeper/zk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ZooKeeper", func() {
	When("A ZooKeeper server is online at localhost", func() {
		It(`Should be able to establish connections
		    and create ephemeral nodes with all 
		    permissions from clients`, func() {
			local := "127.0.0.1"
			conn, _, err := zk.Connect([]string{local}, time.Second)
			Expect(err).To(Equal(nil))

			zkPath := "/zkPath"
			zkFlags := int32(zk.FlagEphemeral)
			data := "zkPathData"
			registeredPath, err := registerUser(conn, zkPath, zkFlags, data)
			Expect(err).To(Equal(nil))
			Expect(registeredPath).To(Equal(zkPath[1:]))

			exists := checkZNode(conn, zkPath)
			Expect(exists).To(Equal(true))

			retrievedData, _ := getZNode(conn, zkPath)
			Expect(retrievedData).To(Equal(data))
		})

		It(`Should be able to establish connections
			and create persistent nodes with all
			permissions from clients`, func() {
			local := "127.0.0.1"
			conn, _, err := zk.Connect([]string{local}, time.Second)
			Expect(err).To(Equal(nil))
		
			zkPath := "/zkPersistent"
			zkFlags := int32(0)
			data := "zkPersistentData"
			registeredPath, err := registerUser(conn, zkPath, zkFlags, data)
			Expect(err).To(Equal(nil))
			Expect(registeredPath).To(Equal(zkPath[1:]))
		})

		It(`Should be able to establish connections
			and retrieve data from persistent nodes`, func() {
			local := "127.0.0.1"
			conn, _, err := zk.Connect([]string{local}, time.Second)
			Expect(err).To(Equal(nil))

			zkPath := "/zkPersistent"
			data := "zkPersistentData"
			exists := checkZNode(conn, zkPath)
			Expect(exists).To(Equal(true))

			retrievedData, _ := getZNode(conn, zkPath)
			Expect(retrievedData).To(Equal(data))
		})
	})
})