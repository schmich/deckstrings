package deckstrings

import (
	"encoding/binary"
  "io"
)

type varintReader struct {
	reader io.ByteReader
}

func (r *varintReader) Read() (uint64, error) {
	return binary.ReadUvarint(r.reader)
}

func (r *varintReader) ReadMany(values []uint64) error {
	for i := 0; i < len(values); i++ {
		if value, err := binary.ReadUvarint(r.reader); err != nil {
			return err
		} else {
			values[i] = value
		}
	}
	return nil
}

type varintWriter struct {
	writer io.Writer
}

func (w *varintWriter) Write(value uint64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, value)
	_, err := w.writer.Write(buf[:n])
	return err
}

func (w *varintWriter) WriteMany(values []uint64) error {
	buf := make([]byte, len(values)*binary.MaxVarintLen64)
	total := 0
	for _, value := range values {
		n := binary.PutUvarint(buf[total:], value)
		total += n
	}
	_, err := w.writer.Write(buf[:total])
	return err
}
