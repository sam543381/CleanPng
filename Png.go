package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	officialSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
)

// Png is a concrete representation of a PNG image
type Png struct {
	signature []byte
	header    header
	data      [][]byte
	palette   palette
}

func (p Png) hasValidSignature() bool {
	return bytes.Compare(officialSignature, p.signature) == 0
}

// DecodePng decodes the given PNG image into a Png struct
func DecodePng(reader io.Reader) (Png, error) {

	png := Png{
		signature: make([]byte, 8),
		data:      make([][]byte, 0),
	}

	reader.Read(png.signature)
	if !png.hasValidSignature() {
		return png, errors.New("Invalid PNG signature")
	}

	chunk, err := decodeChunk(reader)
	if err != nil {
		return png, err
	}

	png.header, err = decodeHeader(bytes.NewBuffer(chunk.data))
	if err != nil {
		log.Fatal(err)
	}

	for chunk, err = decodeChunk(reader); chunk.chunkType != "IEND"; chunk, err = decodeChunk(reader) {
		if err != nil {
			fmt.Printf("Errored chunk : %s", err)
			break
		}

		switch chunk.chunkType {
		case "IDAT":
			png.data = append(png.data, chunk.data)
			break
		case "PLTE":
			png.palette, err = decodePalette(chunk.data)
			if png.header.color == 0x03 { // Palette chunks are mandatory for images with color type 3
				log.Fatal(err)
			} else {
				fmt.Printf("Unable to decode palette chunk: %s\n", err)
			}
		}

	}

	return png, nil
}

// EncodeAndClean encodes a Png image into bytes without its metadata
func (p Png) EncodeAndClean() ([]byte, error) {

	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(officialSignature)

	bytes, err := p.header.encode()
	if err != nil {
		return nil, err
	}
	buffer.Write(bytes)

	if p.header.color == 0x03 { // Palette chunks are mandatory for images with color type 3
		bytes, err = p.palette.encode()
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}

	for i := 0; i < len(p.data); i++ {
		bytes, err = encodeChunk("IDAT", p.data[i])
		if err != nil {
			fmt.Println(err)
		}

		buffer.Write(bytes)
	}

	bytes, err = encodeChunk("IEND", []byte{})
	if err != nil {
		return nil, err
	}
	buffer.Write(bytes)

	return buffer.Bytes(), nil
}

func (p Png) String() string {

	return fmt.Sprintln(
		"Width: ", p.header.width, "px\n",
		"Height: ", p.header.height, "px\n",
		"Depth: ", p.header.depth, "bits\n",
		"Color: ", p.header.color, "\n",
		"Compression: ", p.header.compression, "\n",
		"Filter: ", p.header.filter, "\n",
		"Interlace: ", p.header.interlace, "\n",
		"Signature: ", p.signature, "\n",
		"IDAT chunks: ", len(p.data),
	)

}
