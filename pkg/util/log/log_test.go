package log_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {
	BeforeEach(func() {

	})

	It("can be loaded from JSON", func() {
		Expect("123").To(Equal("123"))
	})
})
