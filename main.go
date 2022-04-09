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
		go scanIndUdpPort(port, &wg)
	}
	wg.Wait()
}

func scanIndUdpPort(port int, wg *sync.WaitGroup) {
	addr, err := net.ResolveUDPAddr(udpProtocol, getPortString(port))
	_, err = net.DialUDP(udpProtocol, nil, addr)
	if noError(err) {
		printPortOpenMsg(udpProtocol, port)
	}
	wg.Done()
}

func ScanTCP(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanIndTcpPort(port, &wg)
	}
	wg.Wait()
}

func scanIndTcpPort(port int, wg *sync.WaitGroup) {
	portStrRep := getPortString(port)
	portScanErr := scanTcpPort(portStrRep)
	if noError(portScanErr) {
		printPortOpenMsg(tcpProtocol, port)
	} else if tooManyConnectionsExists(portScanErr) {
		printError(portScanErr)
	}
	wg.Done()
}

func getPortString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func scanTcpPort(portStrRep string) error {
	_, err := net.Dial(tcpProtocol, portStrRep)
	for tooManyConnectionsExists(err) {
		time.Sleep(time.Nanosecond)
		_, err = net.Dial(tcpProtocol, portStrRep)
	}
	return err
}

func noError(err error) bool {
	return err == nil
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

func printError(err error) {
	fmt.Printf("Error: %v\n", err)
}

func main() {
	ScanUDP(1, 65535)
	ScanTCP(1, 65535)
}
