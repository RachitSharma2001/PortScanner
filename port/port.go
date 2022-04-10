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

type ScanResults struct {
	mu        sync.Mutex
	openPorts []int
}

func (s *ScanResults) addOpenPort(port int) {
	s.mu.Lock()
	s.openPorts = append(s.openPorts, port)
	s.mu.Unlock()
}

func RunWideUDPScan(startPort int, endPort int, hostname string) []int {
	var scanRes ScanResults
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanIndividualUdpPort(port, hostname, &wg, &scanRes)
	}
	wg.Wait()
	return scanRes.openPorts
}

func scanIndividualUdpPort(port int, hostname string, wg *sync.WaitGroup, scanRes *ScanResults) {
	portRep := getPortString(port, hostname)
	addr, err := net.ResolveUDPAddr(udpProtocol, portRep)
	conn, err := net.DialUDP(udpProtocol, nil, addr)
	if errhelp.NoError(err) {
		scanRes.addOpenPort(port)
		conn.Close()
	}
	wg.Done()
}

func RunWideTCPScan(startPort int, endPort int, hostname string) []int {
	var scanRes ScanResults
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go scanIndividualTcpPort(port, hostname, &scanRes, &wg)
	}
	wg.Wait()
	return scanRes.openPorts
}

func scanIndividualTcpPort(port int, hostname string, scanRes *ScanResults, wg *sync.WaitGroup) {
	portStrRep := getPortString(port, hostname)
	portConn, portScanErr := sendRequestToTcpPort(portStrRep)
	if errhelp.NoError(portScanErr) {
		scanRes.addOpenPort(port)
		(*portConn).Close()
	} else if tooManyConnectionsExists(portScanErr) {
		printError(portScanErr)
	}
	wg.Done()
}

func getPortString(port int, hostname string) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func sendRequestToTcpPort(portStrRep string) (*net.Conn, error) {
	conn, err := net.Dial(tcpProtocol, portStrRep)
	for tooManyConnectionsExists(err) {
		time.Sleep(time.Nanosecond)
		conn, err = net.Dial(tcpProtocol, portStrRep)
	}
	return &conn, err
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
