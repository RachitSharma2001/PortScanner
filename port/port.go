package port

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	errhelp "fake.com/PortScanner/errhelp"
)

const (
	tcpProtocol   = "tcp"
	udpProtocol   = "udp"
	noBufferSpace = "system lacked sufficient buffer space"
)

func RunWideUDPScan(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanIndividualUdpPort(port, &wg)
	}
	wg.Wait()
}

func scanIndividualUdpPort(port int, wg *sync.WaitGroup) {
	addr, err := net.ResolveUDPAddr(udpProtocol, getPortString(port))
	_, err = net.DialUDP(udpProtocol, nil, addr)
	if errhelp.NoError(err) {
		printPortOpenMsg(udpProtocol, port)
	}
	wg.Done()
}

func RunWideTCPScan(startPort, endPort int) {
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanIndividualTcpPort(port, &wg)
	}
	wg.Wait()
}

func scanIndividualTcpPort(port int, wg *sync.WaitGroup) {
	portStrRep := getPortString(port)
	portScanErr := sendRequestToTcpPort(portStrRep)
	if errhelp.NoError(portScanErr) {
		printPortOpenMsg(tcpProtocol, port)
	} else if tooManyConnectionsExists(portScanErr) {
		printError(portScanErr)
	}
	wg.Done()
}

func getPortString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func sendRequestToTcpPort(portStrRep string) error {
	_, err := net.Dial(tcpProtocol, portStrRep)
	for tooManyConnectionsExists(err) {
		time.Sleep(time.Nanosecond)
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

func printError(err error) {
	fmt.Printf("Error: %v\n", err)
}
