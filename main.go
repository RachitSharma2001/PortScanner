package main

import (
	"fmt"
	"net"
	"time"
)

func ScanTcp(startPort, endPort int) {
	dialLimit := time.Microsecond
	for port := startPort; port <= endPort; port++ {
		_, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), dialLimit)
		if err == nil {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}

func main() {
	ScanTcp(1, 1024)
}
