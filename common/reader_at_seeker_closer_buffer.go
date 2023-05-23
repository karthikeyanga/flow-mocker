package common

import (
	"bytes"
	"io"
)

type ReaderAtSeekerCloserBuffer struct {
	buf      *bytes.Reader
	isClosed bool
}

func NewReaderAtSeekerCloserBuffer(b []byte) *ReaderAtSeekerCloserBuffer {
	buf := bytes.NewReader(b)
	return &ReaderAtSeekerCloserBuffer{
		buf:      buf,
		isClosed: false,
	}
}

func (b *ReaderAtSeekerCloserBuffer) Read(p []byte) (n int, err error) {
	if b.isClosed {
		return 0, io.EOF
	}
	return b.buf.Read(p)
}

func (b *ReaderAtSeekerCloserBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	if b.isClosed {
		return 0, io.EOF
	}
	return b.buf.ReadAt(p, off)
}

func (b *ReaderAtSeekerCloserBuffer) Seek(offset int64, whence int) (int64, error) {
	if b.isClosed {
		return 0, io.EOF
	}
	return b.buf.Seek(offset, whence)
}

func (b *ReaderAtSeekerCloserBuffer) Close() error {
	b.isClosed = true
	return nil
}
