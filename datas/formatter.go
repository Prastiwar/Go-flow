// Package datas provides a collection of data marshaling and unmarshaling interfaces and functions.
// These interfaces and functions provide a flexible and extensible way to work with structured data in Go,
// whether it is being serialized to a byte slice, written to an IO stream, or both.
// In addition to these interfaces, this package also provides a number of useful functions for working with data in Go,
// including functions for encoding and decoding data using common serialization formats like JSON and XML.
package datas

// ByteFormatter is an interface that combines Marshaler and Unmarshaler into a single interface
// for working with byte slices.
type ByteFormatter interface {
	Marshaler
	Unmarshaler
}

// IOFormatter is an interface that combines WriterMarshaler and ReaderUnmarshaler into a single interface
// for working with IO streams.
type IOFormatter interface {
	WriterMarshaler
	ReaderUnmarshaler
}

// ByteIOFormatter is an interface that combines ByteFormatter and IOFormatter into a single interface
// for working with both byte slices and IO streams.
type ByteIOFormatter interface {
	ByteFormatter
	IOFormatter
}
