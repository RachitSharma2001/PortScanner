package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPortScanner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PortScanner Suite")
}

var _ = Describe("Scanning TCP Ports", func() {
	Context("if we scan an open TCP port", func() {
		It("the scan should be successful", func() {

		})
	})
	Context("if we scan a closed TCP port", func() {
		It("the scan should be unsuccessful", func() {

		})
	})
})
