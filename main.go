package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	tcpProtocol   = "tcp"
	udpProtocol   = "udp"
	noBufferSpace = "system lacked sufficient buffer space"
)

func ScanUDP(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			addr, err := net.ResolveUDPAddr(udpProtocol, getPortString(port))
			_, err = net.DialUDP(udpProtocol, nil, addr)
			if err == nil {
				printPortOpenMsg(udpProtocol, port)
			}
			wg.Done()
		}(port)
	}
	wg.Wait()
}

func ScanTCP(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			portStrRep := getPortString(port)
			portScanErr := scanTcpPort(portStrRep)

			if portScanErr == nil {
				printPortOpenMsg(tcpProtocol, port)
			} else if tooManyConnectionsExists(portScanErr) {
				fmt.Printf("Too many connections error occurred: %v\n", portScanErr)
			}

			wg.Done()
		}(port)
	}
	wg.Wait()
}

func getPortString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func scanTcpPort(portStrRep string) error {
	_, err := net.Dial(tcpProtocol, portStrRep)
	for tooManyConnectionsExists(err) {
		time.Sleep(time.Microsecond)
		_, err = net.Dial(tcpProtocol, portStrRep)
	}
	return err
}

func tooManyConnectionsExists(err error) bool {
	if err == nil {
		return false
	}
	strRepOfErr := err.Error()
	return strings.Contains(strRepOfErr, noBufferSpace)
}

func printPortOpenMsg(protocol string, port int) {
	fmt.Printf("%s port %d is open\n", protocol, port)
}

func main() {
	ScanUDP(1, 65535)
	ScanTCP(1, 65535)
}
