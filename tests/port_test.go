package main_test

import (
	"net"
	"strconv"
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
	host := "localhost"
	portRep := host + ":" + strconv.Itoa(portToListen)
	var listener net.Listener
	Context("if we scan an open TCP port", func() {
		listener, _ = net.Listen("tcp", portRep)
		var scanRes port.ScanResults
		var wg sync.WaitGroup
		wg.Add(1)
		go port.ScanIndividualTcpPort(portToListen, host, &scanRes, &wg)
		wg.Wait()
		openPorts := scanRes.OpenPorts
		It("the scan should detect the open port", func() {
			Expect(openPorts).To(ContainElement(portToListen))
		})
	})
	listener.Close()
	Context("if we scan a closed TCP port", func() {
		var scanRes port.ScanResults
		var wg sync.WaitGroup
		wg.Add(1)
		go port.ScanIndividualTcpPort(portToListen, host, &scanRes, &wg)
		wg.Wait()
		openPorts := scanRes.OpenPorts
		It("the scan should detect no open ports", func() {
			Expect(openPorts).To(BeEmpty())
		})
	})
})
