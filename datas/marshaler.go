package datas

import "io"

var (
	_ Marshaler       = MarshalerFunc(nil)
	_ WriterMarshaler = WriterMarshalerFunc(nil)
)

// Marshaler is an interface for types that can marshal any value to a byte slice.
type Marshaler interface {
	// Marshal returns a byte slice representation of the provided value.
	// An error is returned if the value cannot be marshaled.
	Marshal(v any) ([]byte, error)
}

// WriterMarshalerFunc is a function type that can be used as a Marshaler.
type MarshalerFunc func(v any) ([]byte, error)

func (f MarshalerFunc) Marshal(v any) ([]byte, error) {
	return f(v)
}

// WriterMarshaler is an interface for types that can marshal any value to an io.Writer.
type WriterMarshaler interface {
	// MarshalTo writes a byte slice representation of the provided value to the given io.Writer.
	// An error is returned if the value cannot be marshaled, or if writing to the io.Writer fails.
	MarshalTo(w io.Writer, v any) error
}

// WriterMarshalerFunc is a function type that can be used as a WriterMarshaler.
type WriterMarshalerFunc func(w io.Writer, v any) error

func (f WriterMarshalerFunc) MarshalTo(w io.Writer, v any) error {
	return f(w, v)
}
