package main

import (
	"bytes"
	"errors"
)

type entry struct {
	red   byte
	green byte
	blue  byte
}

func (e entry) encode() ([]byte, error) {

	buffer := bytes.NewBuffer([]byte{})

	buffer.WriteByte(e.red)
	buffer.WriteByte(e.green)
	buffer.WriteByte(e.blue)

	return buffer.Bytes(), nil
}

type palette struct {
	entries []entry
}

func (p palette) encode() ([]byte, error) {

	buffer := bytes.NewBuffer([]byte{})

	for i := 0; i < len(p.entries); i++ {
		bytes, err := p.entries[i].encode()
		if err != nil {
			return buffer.Bytes(), err
		}
		buffer.Write(bytes)
	}

	return encodeChunk("PLTE", buffer.Bytes())

}

func decodePalette(data []byte) (palette, error) {

	palette := palette{
		entries: []entry{},
	}

	if len(data)%3 != 0 {
		return palette, errors.New("Invalid palette length: chunk data length must be a multiple of 3")
	}

	for i := 0; i < len(data)/3; i++ {
		bytes := data[i*3 : (i*3)+3]
		entry, err := decodeEntry(bytes)
		if err != nil {
			return palette, err
		}
		palette.entries = append(palette.entries, entry)
	}

	return palette, nil
}

func decodeEntry(data []byte) (entry, error) {

	return entry{
		red:   data[0],
		green: data[1],
		blue:  data[2],
	}, nil

}
