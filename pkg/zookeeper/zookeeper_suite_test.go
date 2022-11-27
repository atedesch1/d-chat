package zookeeper

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZookeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zookeeper Suite")
}
