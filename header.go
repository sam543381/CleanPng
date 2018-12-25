package main

import (
	"bytes"
	"encoding/binary"
)

type header struct {
	width       uint32
	height      uint32
	depth       uint8
	color       byte
	compression byte
	filter      byte
	interlace   byte
}

func decodeHeader(buffer *bytes.Buffer) (header, error) {

	header := header{}

	bytes := make([]byte, 4)
	buffer.Read(bytes)
	header.width = binary.BigEndian.Uint32(bytes)

	buffer.Read(bytes)
	header.height = binary.BigEndian.Uint32(bytes)

	header.depth, _ = buffer.ReadByte()
	header.color, _ = buffer.ReadByte()
	header.compression, _ = buffer.ReadByte()
	header.filter, _ = buffer.ReadByte()
	header.interlace, _ = buffer.ReadByte()

	return header, nil
}

func (h header) encode() ([]byte, error) {

	buffer := bytes.NewBuffer([]byte{})

	binary.Write(buffer, binary.BigEndian, h.width)
	binary.Write(buffer, binary.BigEndian, h.height)
	buffer.WriteByte(h.depth)
	buffer.WriteByte(h.color)
	buffer.WriteByte(h.compression)
	buffer.WriteByte(h.filter)
	buffer.WriteByte(h.interlace)

	return encodeChunk("IHDR", buffer.Bytes())
}
