package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func ScanTCP(startPort, endPort int) {
	dialLimit := time.Microsecond
	for port := startPort; port <= endPort; port++ {
		_, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), dialLimit)
		if err == nil {
			fmt.Printf("TCP Port %d is open\n", port)
		}
	}
}

func ScanTCPConcurrently(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			portStr := getPortString(port)
			_, err := net.Dial("tcp", portStr)
			if err == nil {
				fmt.Printf("TCP Port %d is open\n", port)
			}
			wg.Done()
		}(port)
	}
	wg.Wait()
}

func getPortString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func main() {
	fmt.Println("Starting Synchronous Scan")
	ScanTCP(1, 1024)
	fmt.Println("Starting Concurrent Scan")
	ScanTCPConcurrently(1, 1024)
}
