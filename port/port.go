package port

import (
	"fmt"
	"net"
	"strings"
	"sync"

	errhelp "fake.com/PortScanner/errhelp"
)

const (
	tcpProtocol   = "tcp"
	udpProtocol   = "udp"
	noBufferSpace = "system lacked sufficient buffer space"
)

type ScanResults struct {
	mu        sync.Mutex
	OpenPorts []int
}

func (s *ScanResults) addOpenPort(port int) {
	s.mu.Lock()
	s.OpenPorts = append(s.OpenPorts, port)
	s.mu.Unlock()
}

func RunWideUDPScan(startPort int, endPort int, hostname string) []int {
	var scanRes ScanResults
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go ScanIndividualUdpPort(port, hostname, &scanRes, &wg)
	}
	wg.Wait()
	return scanRes.OpenPorts
}

func ScanIndividualUdpPort(port int, hostname string, scanRes *ScanResults, wg *sync.WaitGroup) {
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
	portsToScan := make(chan int)
	openPorts := ScanResults{}
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go makeTcpRequests(portsToScan, &openPorts, hostname, &wg)
	}
	addPortsToJobQueue(startPort, endPort, portsToScan)
	wg.Wait()
	return openPorts.OpenPorts
}

func makeTcpRequests(portsToScan chan int, scanRes *ScanResults, hostname string, wg *sync.WaitGroup) {
	for port := range portsToScan {
		conn, err := scanIndividualTcpPort(port, hostname)
		if errhelp.NoError(err) {
			scanRes.addOpenPort(port)
			(*conn).Close()
		} else if tooManyConnectionsExists(err) {
			printError(err)
		}
	}
	wg.Done()
}

func scanIndividualTcpPort(port int, hostname string) (*net.Conn, error) {
	hostnamePortStr := getPortString(port, hostname)
	conn, err := net.Dial(tcpProtocol, hostnamePortStr)
	return &conn, err
}

func tooManyConnectionsExists(err error) bool {
	if err == nil {
		return false
	}
	strRepOfErr := err.Error()
	return strings.Contains(strRepOfErr, noBufferSpace)
}

func addPortsToJobQueue(startPort, endPort int, portsToScan chan int) {
	for port := startPort; port <= endPort; port++ {
		portsToScan <- port
	}
	close(portsToScan)
}

func getPortString(port int, hostname string) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func printError(err error) {
	fmt.Printf("Error: %v\n", err)
}
