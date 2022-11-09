package config

import (
	"io"
	"os"
)

// ReaderDecoder is implemented by any value that has a Decode method.
// The implementation controls how values for corresponding v fields are reader from
// io.Reader and stored in v.
type ReaderDecoder interface {
	Decode(r io.Reader, v any) error
}

type readerProvider struct {
	reader  io.Reader
	decoder ReaderDecoder
}

// NewReaderProvider creates an instance of reader provider with specified reader and decoder.
func NewReaderProvider(r io.Reader, d ReaderDecoder) *readerProvider {
	return &readerProvider{
		reader:  r,
		decoder: d,
	}
}

// Load parses flag definitions from the argument list, which should not include the command name.
// Parsed flag value results are stored in matching v fields. If there is no matching field it
// will be ignored and it's value will not be overridden.
func (p *readerProvider) Load(v any, opts ...LoadOption) error {
	return p.decoder.Decode(p.reader, v)
}

type fileReader struct {
	filename string
}

// NewReader creates an instance of io.Reader reading from file found at filename.
func NewFileReader(filename string) *fileReader {
	return &fileReader{
		filename: filename,
	}
}

// Read opens the named file for reading and reads it up.
func (r *fileReader) Read(p []byte) (n int, err error) {
	f, err := os.Open(r.filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Read(p)
}

// NewReaderProvider returns a new file provider with specified filename and decoder.
func NewFileProvider(filename string, decoder ReaderDecoder) *readerProvider {
	return &readerProvider{
		reader:  NewFileReader(filename),
		decoder: decoder,
	}
}
