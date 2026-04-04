// SPDX-License-Identifier: GPL-2.0-or-later

package capture

import (
	"fmt"
	"io"
)

type ENHCommand byte

const (
	ENHReqInit  ENHCommand = 0x0
	ENHReqSend  ENHCommand = 0x1
	ENHReqStart ENHCommand = 0x2
	ENHReqInfo  ENHCommand = 0x3

	ENHResResetted  ENHCommand = 0x0
	ENHResReceived  ENHCommand = 0x1
	ENHResStarted   ENHCommand = 0x2
	ENHResInfo      ENHCommand = 0x3
	ENHResFailed    ENHCommand = 0xA
	ENHResErrorEBUS ENHCommand = 0xB
	ENHResErrorHost ENHCommand = 0xC
)

const (
	enhByteFlag = byte(0x80)
	enhByteMask = byte(0xC0)
	enhByte1    = byte(0xC0)
	enhByte2    = byte(0x80)
	ebusSync    = byte(0xAA)
)

type ENHFrame struct {
	Command ENHCommand
	Data    byte
}

type ENHParser struct {
	pending bool
	byte1   byte
}

func (parser *ENHParser) Parse(reader io.Reader) (ENHFrame, error) {
	buf := make([]byte, 1)
	for {
		if _, err := io.ReadFull(reader, buf); err != nil {
			return ENHFrame{}, err
		}
		frame, complete, err := parser.Feed(buf[0])
		if err != nil {
			return ENHFrame{}, err
		}
		if complete {
			return frame, nil
		}
	}
}

func (parser *ENHParser) Feed(payloadByte byte) (ENHFrame, bool, error) {
	if !parser.pending {
		if payloadByte&enhByteFlag == 0 {
			return ENHFrame{Command: ENHResReceived, Data: payloadByte}, true, nil
		}
		if payloadByte&enhByteMask == enhByte2 {
			return ENHFrame{}, false, fmt.Errorf("enh malformed frame: orphan second byte 0x%02X", payloadByte)
		}
		parser.pending = true
		parser.byte1 = payloadByte
		return ENHFrame{}, false, nil
	}

	if payloadByte&enhByteMask != enhByte2 {
		parser.pending = false
		return ENHFrame{}, false, fmt.Errorf("enh malformed frame: invalid second byte 0x%02X", payloadByte)
	}

	command := ENHCommand((parser.byte1 >> 2) & 0x0F)
	data := byte(((parser.byte1 & 0x03) << 6) | (payloadByte & 0x3F))
	parser.pending = false
	return ENHFrame{Command: command, Data: data}, true, nil
}

func ENHCommandName(command ENHCommand) string {
	switch command {
	case ENHResResetted:
		return "resetted"
	case ENHResReceived:
		return "received"
	case ENHResStarted:
		return "started"
	case ENHResInfo:
		return "info"
	case ENHResFailed:
		return "failed"
	case ENHResErrorEBUS:
		return "error_ebus"
	case ENHResErrorHost:
		return "error_host"
	default:
		return fmt.Sprintf("command_0x%02X", byte(command))
	}
}
