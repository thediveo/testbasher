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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestCommand", func() {

	It("panics when failing to start a test command", func() {
		Expect(func() { NewTestCommand("") }).To(Panic())
	})

	It("runs a test command", func() {
		c := NewTestCommand("/bin/bash", "-c", `echo "\"go-kay\"" && read`)
		var s string
		c.Decode(&s)
		Expect(s).To(Equal("go-kay"))
		c.Close()
		c.Close()
	})

	It("tells a test command", func() {
		wisdom := "42"
		c := NewTestCommand("/bin/bash", "-c", `read IN && echo "\"${IN}!\"" && read`)
		c.Tell(wisdom)
		var s string
		c.Decode(&s)
		Expect(s).To(Equal(wisdom + "!"))
		c.Close()
	})

	It("ex-terminates a blocking test command", func() {
		c := NewTestCommand("/bin/sleep", "10000001")
		done := make(chan interface{})
		go func() {
			c.Close()
			done <- nil
		}()
		select {
		case <-time.After(5 * time.Second):
			Fail("test command Close() not reacting within time limit")
		case <-done:
		}
	})

	It("returns the commands stderr when decoding fails", func() {
		c := NewTestCommand("/bin/bash", "-c", `
#!/bin/bash
set -e
echo '"42"'
/grmpf
`)
		var s string
		c.Decode(&s)
		Expect(s).To(Equal("42"))
		Expect(func() { c.Decode(&s) }).To(PanicWith(
			MatchRegexp(`(?s)TestCommand\.Decode panicked: EOF
while reading:.*
.*
child process stderr: /bin/bash: line 5: /grmpf: .*`)))
	})

})
