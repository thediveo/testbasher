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
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basher", func() {

	It("creates temporary scripts and cleans up afterwards", func() {
		b := Basher{}
		defer b.Done()

		b.Script("script", `echo "\"$1\"" && echo "\"$script\"" && read`)
		Expect(b.tmpdir).To(BeADirectory())
		Expect(filepath.Join(b.tmpdir, "script.sh")).To(BeARegularFile())

		cmd := b.Start("script", "foo")

		var s string
		cmd.Decode(&s)
		Expect(s).To(Equal("foo"))

		cmd.Decode(&s)
		Expect(s).To(Equal(filepath.Join(b.tmpdir, "script.sh")))

		cmd.Close()
		b.Done()
		Expect(b.tmpdir).ToNot(BeAnExistingFile())
	})

	It("includes common scripts", func() {
		b := Basher{}
		defer b.Done()

		b.Common(`FOOBAR=42`)
		b.Common(`BAZ=12345`)
		b.Script("script", `echo "\"<$FOOBAR><$BAZ>\"" && read`)
		cmd := b.Start("script")
		defer cmd.Close()

		var s string
		cmd.Decode(&s)
		Expect(s).To(Equal("<42><12345>"))
	})

	It("doesn't accept the same script twice", func() {
		b := Basher{}
		defer b.Done()

		Expect(func() { b.Script("foo", "") }).ToNot(Panic())
		Expect(func() { b.Script("foo", "") }).To(Panic())
	})

	It("cannot start an unknown script", func() {
		b := Basher{}
		defer b.Done()

		Expect(func() { b.Start("foo") }).To(Panic())
	})

	It("panics when the filesystem goes wrong", func() {
		b := Basher{}
		Expect(func() { b.init("/nowhere") }).To(Panic())
	})

})
