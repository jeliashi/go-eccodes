package accessor

import "bytes"

type Property int

const (
	GribUserBuffer Property = iota
	GribMyBuffer
)

type GribBuffer struct {
	property     Property
	validity     int
	growable     int
	length       uint32
	ulength      uint32
	ulength_bits uint32
	data         *bytes.Buffer
}

func NewGribBuffer(p Property, data *bytes.Buffer, buflen uint32) *GribBuffer {
	return &GribBuffer{property: GribUserBuffer, length: buflen, ulength: buflen, ulength_bits: 8 * buflen, data: data}
}
