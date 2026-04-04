// SPDX-License-Identifier: GPL-2.0-or-later

package capture

import (
	"encoding/binary"
	"io"
	"time"
)

type PCAPWriter struct {
	writer io.Writer
}

func NewPCAPWriter(writer io.Writer) *PCAPWriter {
	return &PCAPWriter{writer: writer}
}

func (writer *PCAPWriter) WriteHeader() error {
	header := make([]byte, 24)
	binary.LittleEndian.PutUint32(header[0:4], 0xa1b2c3d4)
	binary.LittleEndian.PutUint16(header[4:6], 2)
	binary.LittleEndian.PutUint16(header[6:8], 4)
	binary.LittleEndian.PutUint32(header[8:12], 0)
	binary.LittleEndian.PutUint32(header[12:16], 0)
	binary.LittleEndian.PutUint32(header[16:20], 65535)
	binary.LittleEndian.PutUint32(header[20:24], LinkTypeUSER0)
	_, err := writer.writer.Write(header)
	return err
}

func (writer *PCAPWriter) WritePacket(ts time.Time, payload []byte) error {
	packetHeader := make([]byte, 16)
	binary.LittleEndian.PutUint32(packetHeader[0:4], uint32(ts.Unix()))
	binary.LittleEndian.PutUint32(packetHeader[4:8], uint32(ts.Nanosecond()/1000))
	binary.LittleEndian.PutUint32(packetHeader[8:12], uint32(len(payload)))
	binary.LittleEndian.PutUint32(packetHeader[12:16], uint32(len(payload)))
	if _, err := writer.writer.Write(packetHeader); err != nil {
		return err
	}
	_, err := writer.writer.Write(payload)
	return err
}
