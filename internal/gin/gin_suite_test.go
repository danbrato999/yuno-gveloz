package gin_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gin Suite")
}
