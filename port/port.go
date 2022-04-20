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
		go makeTcpRequests(i+1, portsToScan, &openPorts, hostname, &wg)
	}
	for i := startPort; i <= endPort; i++ {
		portsToScan <- i
	}
	close(portsToScan)
	wg.Wait()
	return openPorts.OpenPorts
}

func makeTcpRequests(workerNum int, portsToScan chan int, scanRes *ScanResults, hostname string, wg *sync.WaitGroup) {
	for port := range portsToScan {
		//fmt.Printf("Worker #%d: port %d\n", workerNum, port)
		portStrRep := getPortString(port, hostname)
		conn, err := net.Dial(tcpProtocol, portStrRep)
		if errhelp.NoError(err) {
			scanRes.addOpenPort(port)
			conn.Close()
		} else if tooManyConnectionsExists(err) {
			fmt.Println("Oh no! - too many connections exist: %v\n", err)
		}

	}
	//	fmt.Printf("Worker #%d: finished\n", workerNum)
	wg.Done()
}

/*func RunWideTCPScan(startPort int, endPort int, hostname string) []int {
	var scanRes ScanResults
	var wg sync.WaitGroup
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go ScanIndividualTcpPort(port, hostname, &scanRes, &wg)
	}
	wg.Wait()
	return scanRes.OpenPorts
}

func ScanIndividualTcpPort(port int, hostname string, scanRes *ScanResults, wg *sync.WaitGroup) {
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

func sendRequestToTcpPort(portStrRep string) (*net.Conn, error) {
	conn, err := net.Dial(tcpProtocol, portStrRep)
	for tooManyConnectionsExists(err) {
		time.Sleep(time.Nanosecond)
		conn, err = net.Dial(tcpProtocol, portStrRep)
	}
	return &conn, err
}*/

func tooManyConnectionsExists(err error) bool {
	if err == nil {
		return false
	}
	strRepOfErr := err.Error()
	return strings.Contains(strRepOfErr, noBufferSpace)
}

func getPortString(port int, hostname string) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func printPortOpenMsg(protocol string, port int) {
	fmt.Printf("%s port %d is open\n", protocol, port)
}

func printError(err error) {
	fmt.Printf("Error: %v\n", err)
}
