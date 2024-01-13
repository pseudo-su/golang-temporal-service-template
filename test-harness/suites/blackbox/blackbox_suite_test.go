package blackbox_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegrationBlackbox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Blackbox Suite")
}
