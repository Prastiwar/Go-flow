package config

import (
	"io"
	"os"
)

type ReaderDecoder interface {
	Decode(r io.Reader, v any) error
}

type readerProvider struct {
	reader  io.Reader
	decoder ReaderDecoder
}

func NewReaderProvider(r io.Reader, d ReaderDecoder) *readerProvider {
	return &readerProvider{
		reader:  r,
		decoder: d,
	}
}

func (p *readerProvider) Load(v any, opts ...LoadOption) error {
	return p.decoder.Decode(p.reader, v)
}

type FileReader struct {
	filename string
}

func NewFileReader(filename string) *FileReader {
	return &FileReader{
		filename: filename,
	}
}

func (r *FileReader) Read(p []byte) (n int, err error) {
	f, err := os.Open(r.filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Read(p)
}

func NewFileProvider(filename string, decoder ReaderDecoder) *readerProvider {
	return &readerProvider{
		reader:  NewFileReader(filename),
		decoder: decoder,
	}
}
