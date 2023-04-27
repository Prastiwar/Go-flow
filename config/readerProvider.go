package config

import (
	"context"
	"io"
	"os"

	"github.com/Prastiwar/Go-flow/datas"
)

type readerProvider struct {
	reader  io.Reader
	decoder datas.ReaderUnmarshaler
}

// NewReaderProvider creates an instance of reader provider with specified reader and decoder.
func NewReaderProvider(r io.Reader, d datas.ReaderUnmarshaler) *readerProvider {
	return &readerProvider{
		reader:  r,
		decoder: d,
	}
}

// Load parses flag definitions from the argument list, which should not include the command name.
// Parsed flag value results are stored in matching v fields. If there is no matching field it
// will be ignored and it's value will not be overridden.
func (p *readerProvider) Load(ctx context.Context, v any, opts ...LoadOption) error {
	return p.decoder.UnmarshalFrom(p.reader, v)
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
func NewFileProvider(filename string, decoder datas.ReaderUnmarshaler) *readerProvider {
	return &readerProvider{
		reader:  NewFileReader(filename),
		decoder: decoder,
	}
}
