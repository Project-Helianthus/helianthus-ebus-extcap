// SPDX-License-Identifier: GPL-2.0-or-later

package capture

const (
	ebusAck  = byte(0x00)
	ebusNack = byte(0xFF)
)

type reconstructPhase uint8

const (
	reconstructPhaseIdle reconstructPhase = iota
	reconstructPhaseRequest
	reconstructPhaseWaitRequestACK
	reconstructPhaseWaitResponse
	reconstructPhaseWaitFinalACK
	reconstructPhaseWaitTerminalSYN
)

type Reconstructor struct {
	phase               reconstructPhase
	requestSegment      []byte
	responseSegment     []byte
	transactionRaw      []byte
	header              EbusTransactionHeader
	responseExpectedLen int
}

func (r *Reconstructor) Feed(frame ENHFrame) (records [][]byte) {
	switch frame.Command {
	case ENHResReceived:
		return r.feedSymbol(frame.Data)
	case ENHResFailed, ENHResErrorEBUS, ENHResErrorHost, ENHResResetted:
		r.reset()
	}
	return nil
}

func (r *Reconstructor) feedSymbol(symbol byte) (records [][]byte) {
	switch r.phase {
	case reconstructPhaseIdle:
		if symbol == ebusSync {
			return nil
		}
		r.phase = reconstructPhaseRequest
		r.requestSegment = append(r.requestSegment[:0], symbol)
		return nil

	case reconstructPhaseRequest:
		if symbol != ebusSync {
			r.requestSegment = append(r.requestSegment, symbol)
			return nil
		}

		header, ok := parseRequestTelegram(r.requestSegment)
		if !ok {
			r.reset()
			return nil
		}
		r.header = header
		switch header.Kind {
		case EbusTransactionBroadcast:
			record := BuildEbusFramePayload(0, false, false, append([]byte(nil), r.requestSegment...))
			r.reset()
			return [][]byte{record}
		case EbusTransactionMasterMaster, EbusTransactionMasterSlave:
			r.phase = reconstructPhaseWaitRequestACK
			return nil
		default:
			r.reset()
			return nil
		}

	case reconstructPhaseWaitRequestACK:
		if !isACKLike(symbol) {
			r.reset()
			return nil
		}

		r.transactionRaw = append(r.transactionRaw[:0], r.requestSegment...)
		r.transactionRaw = append(r.transactionRaw, symbol)
		if symbol == ebusNack || r.header.Kind == EbusTransactionMasterMaster {
			r.phase = reconstructPhaseWaitTerminalSYN
			return nil
		}

		r.phase = reconstructPhaseWaitResponse
		r.responseSegment = r.responseSegment[:0]
		r.responseExpectedLen = 0
		return nil

	case reconstructPhaseWaitResponse:
		// ACK + immediate SYN is treated as an ack-only initiator-responder transaction.
		if len(r.responseSegment) == 0 && symbol == ebusSync {
			record := BuildEbusFramePayload(0, false, false, append([]byte(nil), r.transactionRaw...))
			r.reset()
			return [][]byte{record}
		}
		if symbol == ebusSync {
			r.reset()
			return nil
		}

		r.responseSegment = append(r.responseSegment, symbol)
		if len(r.responseSegment) == 1 {
			r.responseExpectedLen = int(symbol) + 2
		}
		if len(r.responseSegment) < r.responseExpectedLen {
			return nil
		}
		if len(r.responseSegment) > r.responseExpectedLen {
			r.reset()
			return nil
		}

		r.phase = reconstructPhaseWaitFinalACK
		return nil

	case reconstructPhaseWaitFinalACK:
		if !isACKLike(symbol) {
			r.reset()
			return nil
		}
		r.transactionRaw = append(r.transactionRaw, r.responseSegment...)
		r.transactionRaw = append(r.transactionRaw, symbol)
		r.phase = reconstructPhaseWaitTerminalSYN
		return nil

	case reconstructPhaseWaitTerminalSYN:
		if symbol != ebusSync {
			r.reset()
			if symbol != ebusSync {
				r.phase = reconstructPhaseRequest
				r.requestSegment = append(r.requestSegment[:0], symbol)
			}
			return nil
		}
		record := BuildEbusFramePayload(0, false, false, append([]byte(nil), r.transactionRaw...))
		r.reset()
		return [][]byte{record}
	}

	return nil
}

func parseRequestTelegram(raw []byte) (EbusTransactionHeader, bool) {
	header, ok := ParseEbusTransactionHeader(raw)
	if !ok || len(raw) < 6 {
		return EbusTransactionHeader{}, false
	}
	if !isInitiatorCapableAddress(header.QQ) {
		return EbusTransactionHeader{}, false
	}

	requestLen := int(raw[4])
	if len(raw) != 5+requestLen+1 {
		return EbusTransactionHeader{}, false
	}
	return header, true
}

func isACKLike(symbol byte) bool {
	return symbol == ebusAck || symbol == ebusNack
}

func (r *Reconstructor) reset() {
	r.phase = reconstructPhaseIdle
	r.requestSegment = r.requestSegment[:0]
	r.responseSegment = r.responseSegment[:0]
	r.transactionRaw = r.transactionRaw[:0]
	r.header = EbusTransactionHeader{}
	r.responseExpectedLen = 0
}
