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

package testbasher

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("memento Decoder", func() {

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
