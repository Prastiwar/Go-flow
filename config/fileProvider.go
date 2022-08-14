package config

import (
	"io"
	"os"
)

type FileDecoder interface {
	Decode(r io.Reader, v any) error
}

type fileProvider struct {
	filename string
	decoder  FileDecoder
}

func NewFileProvider(filename string, decoder FileDecoder) *fileProvider {
	return &fileProvider{
		filename: filename,
		decoder:  decoder,
	}
}

func (p *fileProvider) Load(v any) error {
	f, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.decoder.Decode(f, v)
}
