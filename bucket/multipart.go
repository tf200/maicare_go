package bucket

import (
	"bytes"
	"mime/multipart"
)

// InMemoryFile implements the multipart.File interface for an in-memory byte slice.
type InMemoryFile struct {
	*bytes.Reader
}

// Close is a no-op that satisfies the io.Closer interface.
func (imf *InMemoryFile) Close() error {
	return nil
}

// This is a compile-time check to ensure InMemoryFile implements multipart.File.
var _ multipart.File = (*InMemoryFile)(nil)
