package testdep

import (
	"net"
	"strconv"
	"strings"
)

func FindFreePort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	l.Close()
	obtainedAddress := l.Addr().String()
	parts := strings.Split(obtainedAddress, ":")
	portString := parts[len(parts)-1]
	port, _ := strconv.Atoi(portString)
	return port, nil
}
