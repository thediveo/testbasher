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
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

// Decoder wraps a json.Decoder using a memento stream reader in order to
// provide actually helpful error messages in case of JSON syntax errors,
// detailing the concrete input stream data where the problem occured.
type Decoder struct {
	// wrapped JSON decoder
	*json.Decoder
	// reminder of the JSON data stream decoded before hitting an error
	m *MementoReader
}

// NewDecoder returns a new memento-enabled JSON decoder, reading from the
// specified reader.
func NewDecoder(r io.Reader) *Decoder {
	m := NewMementoReader(r)
	return &Decoder{
		Decoder: json.NewDecoder(m),
		m:       m,
	}
}

// Decode reads the next JSON-encoded value from its input and stores it in the
// value pointed to by v. In case of a JSON syntax error, the error returned
// contains the decoded input data so far in this Decode() call to give better
// insight of where things went wrong. The error returned wraps the original
// json.SyntaxError object. Other errors get also wrapped with additional
// details about the JSON input read up to now to provide a more meaningful
// context.
func (d *Decoder) Decode(v interface{}) error {
	d.m.Mark(d.Decoder.InputOffset())
	err := d.Decoder.Decode(v)
	if err == nil {
		return nil
	}
	var memento string
	// Get the data the decoder read has read so far in this decoder run.
	if jerr, ok := err.(*json.SyntaxError); ok {
		// If this is a syntax error, then the decoder will tell us at which
		// position in the overall data stream it hit a major road block. Being nice
		// (or not), the position is 1-based, so keep that in mind.
		offset := int(jerr.Offset-d.m.markoffset) - 1
		memento = string(d.m.Memento(d.Decoder.InputOffset() + int64(offset+100)))
		// To provide better context, we then visibly mark the error position in
		// the memento string. Of course, we need to take into account that
		// we're dealing with UTF8 encoded Unicode strings, not wchars or
		// something fixed like that.
		r, rlen := utf8.DecodeRuneInString(memento[offset:])
		memento = fmt.Sprintf("%s►%c◄%s", memento[:offset], r, memento[offset+rlen:])
	} else {
		memento = string(d.m.Memento(d.Decoder.InputOffset() + 100))
	}
	return fmt.Errorf("%w\nwhile reading:\n\t%s", err, memento)
}
