package testdep_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
)

func TestFindFreePort(t *testing.T) {
	port, err := testdep.FindFreePort()
	AssertNoError(t, err)
	if port == 0 {
		fmt.Printf("port is not set")
		t.Fail()
	}
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	AssertNoError(t, err)
	server.Close()
}
