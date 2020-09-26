package tree

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTree(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tree")
}

var _ = Describe("Tree", func() {
	It("origin", func() {
		Expect(origin("website.de/module/test/test.go", "/Users/user/projects", "website.de/module")).
			Should(Equal("/Users/user/projects/test/test.go"))
	})
})
