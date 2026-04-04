// SPDX-License-Identifier: GPL-2.0-or-later

package capture

import (
	"fmt"
	"slices"
)

type StreamKind string

const (
	StreamENSEvents  StreamKind = "ens-events"
	StreamEbusFrames StreamKind = "ebus-frames"
	StreamBoth       StreamKind = "both"
)

type Direction string

const (
	DirectionUnknown  Direction = "unknown"
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
)

type Record struct {
	RecordVersion  uint8
	StreamKind     StreamKind
	Direction      Direction
	TransportClass string
	EndpointKind   string
	PB             uint8
	SB             uint8
	Family         string
	FrameType      string
	Outcome        string
	RawPayload     []byte
}

type Filter struct {
	SourceAddress *uint8
	TargetAddress *uint8
	Opcode        *uint16
	Family        string
}

var SupportedSemanticEbusOpcodes = []uint16{
	0x0304, 0x0305, 0x0306, 0x0307, 0x0308, 0x0310,
	0x0500, 0x0501, 0x0502, 0x0503, 0x0504, 0x0506, 0x0507, 0x0508, 0x0509, 0x050A, 0x050B, 0x050C, 0x050D,
	0x0700, 0x0701, 0x0702, 0x0703, 0x0704, 0x07FE, 0x07FF,
	0x0800, 0x0801, 0x0802, 0x0803, 0x0804,
	0x0900, 0x0901, 0x0902, 0x0903,
	0x0F01, 0x0F02, 0x0F03,
	0xFE01,
	0xFF00, 0xFF01, 0xFF02, 0xFF03, 0xFF04, 0xFF05, 0xFF06,
}

var SupportedSemanticVaillantOpcodes = []uint16{
	0xB504, 0xB505, 0xB506, 0xB509, 0xB510, 0xB511,
	0xB512, 0xB516, 0xB51A, 0xB524, 0xB555,
}

var semanticLabels = map[uint16]string{
	0x0304: "eBUS semantic 0x0304",
	0x0305: "eBUS semantic 0x0305",
	0x0306: "eBUS semantic 0x0306",
	0x0307: "eBUS semantic 0x0307",
	0x0308: "eBUS semantic 0x0308",
	0x0310: "eBUS semantic 0x0310",
	0x0500: "eBUS semantic 0x0500",
	0x0501: "eBUS semantic 0x0501",
	0x0502: "eBUS semantic 0x0502",
	0x0503: "eBUS semantic 0x0503",
	0x0504: "eBUS semantic 0x0504",
	0x0506: "eBUS semantic 0x0506",
	0x0507: "eBUS semantic 0x0507",
	0x0508: "eBUS semantic 0x0508",
	0x0509: "eBUS semantic 0x0509",
	0x050A: "eBUS semantic 0x050A",
	0x050B: "eBUS semantic 0x050B",
	0x050C: "eBUS semantic 0x050C",
	0x050D: "eBUS semantic 0x050D",
	0x0700: "eBUS semantic 0x0700",
	0x0701: "eBUS semantic 0x0701",
	0x0702: "eBUS semantic 0x0702",
	0x0703: "eBUS semantic 0x0703",
	0x0704: "eBUS semantic 0x0704",
	0x07FE: "eBUS semantic 0x07FE",
	0x07FF: "eBUS semantic 0x07FF",
	0x0800: "eBUS semantic 0x0800",
	0x0801: "eBUS semantic 0x0801",
	0x0802: "eBUS semantic 0x0802",
	0x0803: "eBUS semantic 0x0803",
	0x0804: "eBUS semantic 0x0804",
	0x0900: "eBUS semantic 0x0900",
	0x0901: "eBUS semantic 0x0901",
	0x0902: "eBUS semantic 0x0902",
	0x0903: "eBUS semantic 0x0903",
	0x0F01: "eBUS semantic 0x0F01",
	0x0F02: "eBUS semantic 0x0F02",
	0x0F03: "eBUS semantic 0x0F03",
	0xFE01: "eBUS semantic 0xFE01",
	0xFF00: "eBUS semantic 0xFF00",
	0xFF01: "eBUS semantic 0xFF01",
	0xFF02: "eBUS semantic 0xFF02",
	0xFF03: "eBUS semantic 0xFF03",
	0xFF04: "eBUS semantic 0xFF04",
	0xFF05: "eBUS semantic 0xFF05",
	0xFF06: "eBUS semantic 0xFF06",
	0xB504: "Vaillant B504",
	0xB505: "Vaillant B505",
	0xB506: "Vaillant B506",
	0xB509: "Vaillant B509 boiler registers",
	0xB510: "Vaillant B510",
	0xB511: "Vaillant B511",
	0xB512: "Vaillant B512",
	0xB516: "Vaillant B516 energy",
	0xB51A: "Vaillant B51A",
	0xB524: "Vaillant B524 extended registers",
	0xB555: "Vaillant B555 timer protocol",
}

func OpcodeFromPBSB(pb, sb uint8) uint16 {
	return uint16(pb)<<8 | uint16(sb)
}

func IsSupportedSemanticOpcode(opcode uint16) bool {
	return slices.Contains(SupportedSemanticEbusOpcodes, opcode) ||
		slices.Contains(SupportedSemanticVaillantOpcodes, opcode)
}

func SemanticLabel(opcode uint16) string {
	if label, ok := semanticLabels[opcode]; ok {
		return label
	}
	return fmt.Sprintf("Raw opcode 0x%04X", opcode)
}
