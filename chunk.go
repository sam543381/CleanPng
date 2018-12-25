package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type chunk struct {
	length    uint32
	chunkType string
	data      []byte
	crc       uint32
}

func (c chunk) hasValidCrc() bool {

	buffer := bytes.NewBufferString(c.chunkType)
	buffer.Write(c.data)

	return c.crc == crc32.ChecksumIEEE(buffer.Bytes())
}

func decodeChunk(reader io.Reader) (chunk, error) {

	chunk := chunk{}

	bytes := make([]byte, 4)
	reader.Read(bytes)
	chunk.length = binary.BigEndian.Uint32(bytes)

	reader.Read(bytes)
	chunk.chunkType = string(bytes)

	chunk.data = make([]byte, chunk.length)
	reader.Read(chunk.data)

	reader.Read(bytes)
	chunk.crc = binary.BigEndian.Uint32(bytes)

	if !chunk.hasValidCrc() {
		return chunk, errors.New("Chunk has an invalid crc")
	}

	return chunk, nil
}

func encodeChunk(chunkType string, data []byte) ([]byte, error) {

	buffer := bytes.NewBuffer([]byte{})

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(data)))
	buffer.Write(lengthBytes)

	binary.Write(buffer, binary.BigEndian, len(data))
	buffer.WriteString(chunkType)
	buffer.Write(data)

	crcBuffer := bytes.NewBufferString(chunkType)
	crcBuffer.Write(data)
	binary.Write(buffer, binary.BigEndian, crc32.ChecksumIEEE(crcBuffer.Bytes()))

	return buffer.Bytes(), nil
}
