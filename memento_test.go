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

// +build go1.14

package testbasher

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memento", func() {

	It("reads and remembers", func() {
		m := NewMementoReader(strings.NewReader("1234567890"))
		Expect(m).NotTo(BeNil())
		Expect(m.Memento(100)).To(BeEmpty())

		b := make([]byte, 1)
		n, err := m.Read(b)
		Expect(n).To(Equal(1))
		Expect(err).NotTo(HaveOccurred())
		Expect(m.Memento(100)).To(Equal([]byte("1")))

		b = make([]byte, 4)
		n, err = m.Read(b)
		Expect(n).To(Equal(4))
		Expect(err).NotTo(HaveOccurred())
		Expect(m.Memento(100)).To(Equal([]byte("12345")))

		m.Mark(3)
		Expect(m.Memento(5)).To(Equal([]byte("45")))

		m.Mark(5)
		b = make([]byte, 100)
		n, err = m.Read(b)
		Expect(n).To(Equal(5))
		Expect(m.Memento(100)).To(Equal([]byte("67890")))

		m.Mark(10)
		b = make([]byte, 1)
		n, err = m.Read(b)
		Expect(n).To(BeZero())
		Expect(err).To(Equal(io.EOF))
	})

})
