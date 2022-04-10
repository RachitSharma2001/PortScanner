package main_test

import (
	"fmt"
	"net"
	"sync"
	"testing"

	port "fake.com/PortScanner/port"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPortScanner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PortScanner Suite")
}

var _ = Describe("Scanning TCP Ports", func() {
	portToListen := 3000
	hostname := "localhost"
	var listener net.Listener
	Context("if we scan an open TCP port", func() {
		setUpServerAtTcpPort(&listener, portToListen, hostname)
		openPorts := getOpenTcpPorts(portToListen, hostname)
		It("the scan should detect the open port", func() {
			Expect(openPorts).To(ContainElement(portToListen))
		})
	})
	listener.Close()
	Context("if we scan a closed TCP port", func() {
		openPorts := getOpenTcpPorts(portToListen, hostname)
		It("the scan should detect no open ports", func() {
			Expect(openPorts).To(BeEmpty())
		})
	})
})

var _ = Describe("Scanning UDP Ports", func() {
	portToListen := 3000
	hostname := "localhost"
	var listener net.PacketConn
	Context("if we scan an open UDP port", func() {
		setUpServerAtUdpPort(&listener, portToListen, hostname)
		openPorts := getOpenUdpPorts(portToListen, hostname)
		It("the scan should detect the open port", func() {
			Expect(openPorts).To(ContainElement(portToListen))
		})
	})
	listener.Close()
})

func setUpServerAtTcpPort(listener *net.Listener, portToListen int, hostname string) {
	portRep := getPortRep(portToListen, hostname)
	*listener, _ = net.Listen("tcp", portRep)
}

func getPortRep(port int, hostname string) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}

func getOpenTcpPorts(portToListen int, hostname string) []int {
	var scanRes port.ScanResults
	var wg sync.WaitGroup
	wg.Add(1)
	go port.ScanIndividualTcpPort(portToListen, hostname, &scanRes, &wg)
	wg.Wait()
	return scanRes.OpenPorts
}

func setUpServerAtUdpPort(listener *net.PacketConn, portToListen int, hostname string) {
	portRep := getPortRep(portToListen, hostname)
	*listener, _ = net.ListenPacket("udp", portRep)
}

func getOpenUdpPorts(portToListen int, hostname string) []int {
	var scanRes port.ScanResults
	var wg sync.WaitGroup
	wg.Add(1)
	go port.ScanIndividualUdpPort(portToListen, hostname, &scanRes, &wg)
	wg.Wait()
	return scanRes.OpenPorts
}
