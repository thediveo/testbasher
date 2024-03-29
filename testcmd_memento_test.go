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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCommand with memento", func() {

	It("panics on undecodable test command output with helpful details", func() {
		c := NewTestCommand("/bin/bash", "-c", `echo \"foo && read`)
		var s string
		Expect(func() { c.Decode(&s) }).To(PanicWith(MatchRegexp(
			"(?s)TestCommand\\.Decode panicked: invalid character '\\\\n' in string literal\nwhile reading:\n\t\"foo►\n◄\nchild process stderr:.*")))
	})

})
