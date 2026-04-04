// SPDX-License-Identifier: GPL-2.0-or-later

package capture

import (
	"encoding/binary"
	"fmt"
)

const (
	RecordKindENSEvent  byte = 1
	RecordKindEbusFrame byte = 2
)

const (
	recordFlagHasSource      byte = 1 << 0
	recordFlagRequestLike    byte = 1 << 1
	recordFlagSyncTerminated byte = 1 << 2
)

const (
	RecordVersion byte   = 1
	LinkTypeUSER0 uint32 = 147
)

func BuildENSEventPayload(frame ENHFrame) []byte {
	raw := []byte{frame.Data}
	return buildRecordPayload(RecordKindENSEvent, byte(frame.Command), 0, 0, 0, 0, raw)
}

func BuildEbusFramePayload(source byte, hasSource bool, requestLike bool, raw []byte) []byte {
	flags := recordFlagSyncTerminated
	if hasSource {
		flags |= recordFlagHasSource
	}
	if requestLike {
		flags |= recordFlagRequestLike
	}

	pb, sb := detectPBSB(raw, hasSource)
	return buildRecordPayload(RecordKindEbusFrame, 0, flags, source, pb, sb, raw)
}

func buildRecordPayload(kind byte, command byte, flags byte, source byte, pb byte, sb byte, raw []byte) []byte {
	payload := make([]byte, 14+len(raw))
	copy(payload[:4], []byte("HLTH"))
	payload[4] = RecordVersion
	payload[5] = kind
	payload[6] = command
	payload[7] = flags
	payload[8] = source
	payload[9] = pb
	payload[10] = sb
	payload[11] = 0
	binary.LittleEndian.PutUint16(payload[12:14], uint16(len(raw)))
	copy(payload[14:], raw)
	return payload
}

func detectPBSB(raw []byte, hasSource bool) (byte, byte) {
	_ = hasSource
	if header, ok := ParseEbusTransactionHeader(raw); ok {
		return header.PB, header.SB
	}
	return 0, 0
}

func ParseRecordPayload(data []byte) (kind byte, command byte, flags byte, source byte, pb byte, sb byte, raw []byte, err error) {
	if len(data) < 14 {
		return 0, 0, 0, 0, 0, 0, nil, fmt.Errorf("record too short: %d", len(data))
	}
	if string(data[:4]) != "HLTH" {
		return 0, 0, 0, 0, 0, 0, nil, fmt.Errorf("record magic mismatch")
	}
	if data[4] != RecordVersion {
		return 0, 0, 0, 0, 0, 0, nil, fmt.Errorf("record version mismatch: %d", data[4])
	}

	rawLen := int(binary.LittleEndian.Uint16(data[12:14]))
	if len(data) < 14+rawLen {
		return 0, 0, 0, 0, 0, 0, nil, fmt.Errorf("record payload truncated: have %d want %d", len(data), 14+rawLen)
	}

	return data[5], data[6], data[7], data[8], data[9], data[10], data[14 : 14+rawLen], nil
}
