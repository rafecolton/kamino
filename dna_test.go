package kamino_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/modcloth/kamino"
)

var _ = Describe("NewGenome()", func() {

	var opts map[string]string

	BeforeEach(func() {
		opts = map[string]string{
			"depth":   "50",
			"token":   "abc123",
			"account": "modcloth",
			"repo":    "kamino",
			"cache":   "",
			"ref":     "123",
		}
	})

	It("assigns token from the provided options", func() {
		subject, _ := NewGenome(opts)

		Expect(subject.APIToken).To(Equal("abc123"))
	})

	Context("with a non-integer depth", func() {
		It("returns an error", func() {
			opts["depth"] = "foo"
			subject, err := NewGenome(opts)

			Expect(subject).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no account specified", func() {
		It("returns an error", func() {
			opts["account"] = ""
			subject, err := NewGenome(opts)

			Expect(subject).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no repo specified", func() {
		It("returns an error", func() {
			opts["repo"] = ""
			subject, err := NewGenome(opts)

			Expect(subject).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with an empty cache option", func() {
		It("defaults the cache option to `no`", func() {
			subject, _ := NewGenome(opts)

			Expect(subject.UseCache).To(Equal("no"))
		})
	})

	Context("with an invalid cache option", func() {
		It("returns an error", func() {
			opts["cache"] = "foo"
			subject, err := NewGenome(opts)

			Expect(subject).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})

	Context("with no ref specified", func() {
		It("returns an error", func() {
			opts["ref"] = ""
			subject, err := NewGenome(opts)

			Expect(subject).To(BeNil())
			Expect(err).ToNot(BeNil())
		})
	})
})
