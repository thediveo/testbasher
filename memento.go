package testbasher

import "io"

// MementoReader is an io.Reader wrapping another io.Reader and remembering what
// has been read so far, until it is allowed to forget by starting a new memory
// cycle using Mark().
type MementoReader struct {
	reader     io.Reader // wrapped/source reader
	memento    []byte    // stream data buffered so far
	pos        int       // current reading position in memento
	markoffset int64     // offset of first memento byte since reading started
}

// NewMementoReader returns a new MemontoReader wrapping the specified
// io.Reader.
func NewMementoReader(r io.Reader) *MementoReader {
	return &MementoReader{
		reader:  r,
		memento: []byte{},
	}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0
// <= n <= len(p)) and any error encountered. Even if Read returns n < len(p),
// it may use all of p as scratch space during the call. If some data is
// available but not len(p) bytes, Read conventionally returns what is available
// instead of waiting for more.
func (m *MementoReader) Read(p []byte) (n int, err error) {
	toread := len(p)
	if missing := toread - (len(m.memento) - m.pos); missing > 0 {
		// We need to read in more data from the wrapper reader in order to
		// fulfill the caller's Read() request. But first, let's allocate enough
		// additional room in our memento.
		newbuf := make([]byte, len(m.memento)+missing)
		copy(newbuf, m.memento)
		m.memento = newbuf
		// Now try to get the requested amount of data.
		read := 0
		for read < toread {
			var n int
			n, err = m.reader.Read(m.memento[m.pos+read : m.pos+toread])
			read += n
			if err != nil || n < toread-read {
				break
			}
		}
		if read < toread {
			m.memento = m.memento[:m.pos+read]
		}
		// Whatever amount of data we could read, return it, updating our
		// internal position, but still remembering all we've read so far.
		copy(p, m.memento)
		m.pos += read
		return read, err
	}
	m.pos += toread
	return toread, nil
}

// Mark sets the beginning of the memento, forgetting the previous memento. All
// data read from now on will be remembered until the next Mark().
func (m *MementoReader) Mark(offset int64) {
	if trash := int(offset - m.markoffset); trash > 0 {
		copy(m.memento, m.memento[trash:])
		m.memento = m.memento[:len(m.memento)-trash]
		m.markoffset = offset
		m.pos = 0
	}
}

// Memento returns the data read since the last Mark(). This does not reset the
// memento.
func (m *MementoReader) Memento(offset int64) []byte {
	if offset-m.markoffset > int64(len(m.memento)) {
		return m.memento
	}
	return m.memento[:offset-m.markoffset]
}
