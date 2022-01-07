// Copyright 2020 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

//go:build go1.14
// +build go1.14

package testbasher

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// This is going to be a "hardcore" test where we use a string reader that
// never returns enough string even if it could, but at most a single rune.
type ThickStringReader struct {
	*strings.Reader
}

func NewReader(s string) *ThickStringReader {
	return &ThickStringReader{strings.NewReader(s)}
}

func (r *ThickStringReader) Read(b []byte) (n int, err error) {
	bminor := b[0:1]
	return r.Reader.Read(bminor)
}

var _ = Describe("Memento", func() {

	It("reads and remembers", func() {
		m := NewMementoReader(NewReader("1234567890"))
		Expect(m).NotTo(BeNil())
		Expect(m.Memento(100)).To(BeEmpty())

		b := make([]byte, 1)
		n, err := m.Read(b)
		Expect(n).To(Equal(1))
		Expect(err).To(Succeed())
		Expect(m.Memento(100)).To(Equal([]byte("1")))

		b = make([]byte, 4)
		for _, r := range "2345" {
			n, err = m.Read(b)
			Expect(n).To(Equal(1))
			Expect(err).To(Succeed())
			Expect(string(b[0])).To(Equal(string(r)))
		}
		Expect(m.Memento(100)).To(Equal([]byte("12345")))

		m.Mark(3)
		Expect(m.Memento(5)).To(Equal([]byte("45")))

		m.Mark(5)
		b = make([]byte, 100)
		for _, r := range "67890x" {
			n, err = m.Read(b)
			if r == 'x' {
				Expect(n).To(BeZero())
				Expect(err).To(Equal(io.EOF))
				break
			}
			Expect(n).To(Equal(1))
			Expect(err).To(Succeed())
			Expect(string(b[0])).To(Equal(string(r)))
		}
		Expect(m.Memento(100)).To(Equal([]byte("67890")))

		m.Mark(10)
		b = make([]byte, 1)
		n, err = m.Read(b)
		Expect(n).To(BeZero())
		Expect(err).To(Equal(io.EOF))

		b = make([]byte, 0)
		n, err = m.Read(b)
		Expect(n).To(BeZero())
		Expect(err).To(Succeed())
	})

})
