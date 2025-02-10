package stores_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStores(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stores Suite")
}
