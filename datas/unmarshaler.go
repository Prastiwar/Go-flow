package datas

import "io"

var (
	_ Unmarshaler       = UnmarshalerFunc(nil)
	_ ReaderUnmarshaler = ReaderUnmarshalerFunc(nil)
)

// Unmarshaler is an interface for types that can unmarshal a byte slice to any value.
type Unmarshaler interface {
	// Unmarshal populates the given value with data from the byte slice.
	// An error is returned if the value cannot be unmarshaled.
	Unmarshal(data []byte, v any) error
}

// UnmarshalerFunc is a function type that can be used as a Unmarshaler.
type UnmarshalerFunc func(data []byte, v any) error

func (f UnmarshalerFunc) Unmarshal(data []byte, v any) error {
	return f(data, v)
}

// ReaderUnmarshaler is an interface for types that can unmarshal data from an io.Reader into a value.
type ReaderUnmarshaler interface {
	// UnmarshalFrom reads data from the given io.Reader and populates the given value with it.
	// An error is returned if the data cannot be unmarshaled, or if reading from the io.Reader fails.
	UnmarshalFrom(r io.Reader, v any) error
}

// ReaderUnmarshalerFunc is a function type that can be used as a ReaderUnmarshaler.
type ReaderUnmarshalerFunc func(r io.Reader, v any) error

func (f ReaderUnmarshalerFunc) UnmarshalFrom(r io.Reader, v any) error {
	return f(r, v)
}
