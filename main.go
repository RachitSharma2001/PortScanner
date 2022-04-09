package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	noBufferSpace = "system lacked sufficient buffer space"
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
			portStrRep := getPortString(port)
			_, err := net.Dial("tcp", portStrRep)

			for tooManyConnections(err) {
				time.Sleep(time.Microsecond)
				_, err = net.Dial("tcp", portStrRep)
			}

			if err == nil {
				fmt.Printf("TCP Port %d is open\n", port)
			} else if tooManyConnections(err) {
				fmt.Printf("Too many connections error occurred: %v\n", err)
			}

			wg.Done()
		}(port)
	}
	wg.Wait()
}

func getPortString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func tooManyConnections(err error) bool {
	if err == nil {
		return false
	}
	strRepOfErr := err.Error()
	return strings.Contains(strRepOfErr, noBufferSpace)
}

func main() {
	ScanTCPConcurrently(1, 65535)
}
