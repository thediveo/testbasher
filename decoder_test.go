package testbasher

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {

	It("decodes", func() {
		d := NewDecoder(strings.NewReader("42\n\"abc\""))
		var i int
		Expect(d.Decode(&i)).NotTo(HaveOccurred())
		Expect(i).To(Equal(42))
		var s string
		Expect(d.Decode(&s)).NotTo(HaveOccurred())
		Expect(s).To(Equal("abc"))
	})

	It("gives useful error messages", func() {
		d := NewDecoder(strings.NewReader("42\n{\"foo\":\"bar\", foobar}"))
		var i int
		Expect(d.Decode(&i)).NotTo(HaveOccurred())
		Expect(i).To(Equal(42))
		var s string
		Expect(d.Decode(&s)).To(MatchError(MatchRegexp(
			`invalid character 'f' looking for beginning .+
while reading:
\t
{"foo":"bar", ►f◄oobar}`)))
	})

})
