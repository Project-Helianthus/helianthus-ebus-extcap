// SPDX-License-Identifier: GPL-2.0-or-later

package capture

const ebusBroadcast = byte(0xFE)

type EbusTransactionKind string

const (
	EbusTransactionUnknown      EbusTransactionKind = "unknown"
	EbusTransactionBroadcast    EbusTransactionKind = "broadcast"
	EbusTransactionMasterMaster EbusTransactionKind = "master_master"
	EbusTransactionMasterSlave  EbusTransactionKind = "master_slave"
)

type EbusTransactionHeader struct {
	QQ   byte
	ZZ   byte
	PB   byte
	SB   byte
	Kind EbusTransactionKind
}

func ParseEbusTransactionHeader(raw []byte) (EbusTransactionHeader, bool) {
	if len(raw) < 4 {
		return EbusTransactionHeader{}, false
	}

	header := EbusTransactionHeader{
		QQ: raw[0],
		ZZ: raw[1],
		PB: raw[2],
		SB: raw[3],
	}

	switch {
	case header.ZZ == ebusBroadcast:
		header.Kind = EbusTransactionBroadcast
	case isInitiatorCapableAddress(header.ZZ):
		header.Kind = EbusTransactionMasterMaster
	default:
		header.Kind = EbusTransactionMasterSlave
	}

	return header, true
}

func isInitiatorCapableAddress(addr byte) bool {
	return initiatorPartIndex(addr&0x0F) > 0 && initiatorPartIndex((addr&0xF0)>>4) > 0
}

func initiatorPartIndex(bits byte) byte {
	switch bits {
	case 0x0:
		return 1
	case 0x1:
		return 2
	case 0x3:
		return 3
	case 0x7:
		return 4
	case 0xF:
		return 5
	default:
		return 0
	}
}
