package capture

import (
	"bytes"
	"testing"
	"time"
)

func TestBuildAndParseENSEventRecord(t *testing.T) {
	payload := BuildENSEventPayload(ENHFrame{Command: ENHResStarted, Data: 0x31})
	kind, command, _, source, pb, sb, raw, err := ParseRecordPayload(payload)
	if err != nil {
		t.Fatalf("ParseRecordPayload error = %v", err)
	}
	if kind != RecordKindENSEvent || command != byte(ENHResStarted) {
		t.Fatalf("kind=%d command=%d", kind, command)
	}
	if source != 0 || pb != 0 || sb != 0 {
		t.Fatalf("unexpected metadata source=%d pb=%d sb=%d", source, pb, sb)
	}
	if len(raw) != 1 || raw[0] != 0x31 {
		t.Fatalf("raw = %#v", raw)
	}
}

func TestReconstructorBuildsMasterMasterRecord(t *testing.T) {
	var r Reconstructor
	out := feedENHSequence(&r,
		0x71, 0x10, 0xB5, 0x24, 0x01, 0x02, 0xD1,
		ebusSync, ebusAck, ebusSync,
	)
	if len(out) != 1 {
		t.Fatalf("len(out) = %d; want 1", len(out))
	}
	kind, _, flags, source, pb, sb, raw, err := ParseRecordPayload(out[0])
	if err != nil {
		t.Fatalf("ParseRecordPayload error = %v", err)
	}
	if kind != RecordKindEbusFrame {
		t.Fatalf("kind = %d", kind)
	}
	if source != 0x00 || pb != 0xB5 || sb != 0x24 {
		t.Fatalf("source=%#x pb=%#x sb=%#x raw=%x", source, pb, sb, raw)
	}
	if flags&recordFlagHasSource != 0 || flags&recordFlagRequestLike != 0 {
		t.Fatalf("flags = %#x", flags)
	}
	if got, want := raw, []byte{0x71, 0x10, 0xB5, 0x24, 0x01, 0x02, 0xD1, 0x00}; !bytes.Equal(got, want) {
		t.Fatalf("raw=%x want=%x", got, want)
	}
}

func TestReconstructorBuildsMasterSlaveAckOnlyRecordWithB509Opcode(t *testing.T) {
	var r Reconstructor
	out := feedENHSequence(&r,
		0x71, 0x08, 0xB5, 0x09, 0x03, 0x0D, 0xAB, 0x00, 0xD1,
		ebusSync, ebusAck, ebusSync,
	)
	if len(out) != 1 {
		t.Fatalf("len(out) = %d; want 1", len(out))
	}

	kind, _, flags, source, pb, sb, raw, err := ParseRecordPayload(out[0])
	if err != nil {
		t.Fatalf("ParseRecordPayload error = %v", err)
	}
	if kind != RecordKindEbusFrame {
		t.Fatalf("kind = %d", kind)
	}
	if source != 0x00 {
		t.Fatalf("source=%#x; want 0", source)
	}
	if pb != 0xB5 || sb != 0x09 {
		t.Fatalf("pb=%#x sb=%#x raw=%x", pb, sb, raw)
	}
	if flags&recordFlagHasSource != 0 || flags&recordFlagRequestLike != 0 {
		t.Fatalf("flags = %#x", flags)
	}
	if got, want := raw, []byte{0x71, 0x08, 0xB5, 0x09, 0x03, 0x0D, 0xAB, 0x00, 0xD1, 0x00}; !bytes.Equal(got, want) {
		t.Fatalf("raw=%x want=%x", got, want)
	}
}

func TestReconstructorBuildsMasterSlaveRecordWithResponseSegment(t *testing.T) {
	var r Reconstructor
	out := feedENHSequence(&r,
		0x71, 0x15, 0xB5, 0x24, 0x01, 0x02, 0xD1,
		ebusSync, ebusAck,
		0x02, 0x11, 0x22, 0xE3,
		ebusAck, ebusSync,
	)
	if len(out) != 1 {
		t.Fatalf("len(out) = %d; want 1", len(out))
	}

	_, _, _, _, pb, sb, raw, err := ParseRecordPayload(out[0])
	if err != nil {
		t.Fatalf("ParseRecordPayload error = %v", err)
	}
	if pb != 0xB5 || sb != 0x24 {
		t.Fatalf("pb=%#x sb=%#x raw=%x", pb, sb, raw)
	}
	if got, want := raw, []byte{0x71, 0x15, 0xB5, 0x24, 0x01, 0x02, 0xD1, 0x00, 0x02, 0x11, 0x22, 0xE3, 0x00}; !bytes.Equal(got, want) {
		t.Fatalf("raw=%x want=%x", got, want)
	}
}

func TestReconstructorBuildsBroadcastRecord(t *testing.T) {
	var r Reconstructor
	out := feedENHSequence(&r,
		0x31, 0xFE, 0xB5, 0x24, 0x01, 0x00, 0xD1,
		ebusSync,
	)
	if len(out) != 1 {
		t.Fatalf("len(out) = %d; want 1", len(out))
	}

	_, _, _, _, pb, sb, raw, err := ParseRecordPayload(out[0])
	if err != nil {
		t.Fatalf("ParseRecordPayload error = %v", err)
	}
	if pb != 0xB5 || sb != 0x24 {
		t.Fatalf("pb=%#x sb=%#x raw=%x", pb, sb, raw)
	}
	if got, want := raw, []byte{0x31, 0xFE, 0xB5, 0x24, 0x01, 0x00, 0xD1}; !bytes.Equal(got, want) {
		t.Fatalf("raw=%x want=%x", got, want)
	}
}

func TestParseEbusTransactionHeaderTreatsSlaveTargetAsMasterSlave(t *testing.T) {
	header, ok := ParseEbusTransactionHeader([]byte{0x71, 0x08, 0xB5, 0x09, 0x03, 0x0D, 0xAB, 0x00, 0xD1})
	if !ok {
		t.Fatal("ParseEbusTransactionHeader() ok = false")
	}
	if header.QQ != 0x71 || header.ZZ != 0x08 || header.PB != 0xB5 || header.SB != 0x09 {
		t.Fatalf("header = %#v", header)
	}
	if header.Kind != EbusTransactionMasterSlave {
		t.Fatalf("header.Kind = %q; want %q", header.Kind, EbusTransactionMasterSlave)
	}
}

func TestParseEbusTransactionHeaderTreatsBroadcastAsBroadcast(t *testing.T) {
	header, ok := ParseEbusTransactionHeader([]byte{0x31, 0xFE, 0xB5, 0x24, 0x01})
	if !ok {
		t.Fatal("ParseEbusTransactionHeader() ok = false")
	}
	if header.Kind != EbusTransactionBroadcast {
		t.Fatalf("header.Kind = %q; want %q", header.Kind, EbusTransactionBroadcast)
	}
}

func TestPCAPWriterWritesHeaderAndPacket(t *testing.T) {
	var buf bytes.Buffer
	writer := NewPCAPWriter(&buf)
	if err := writer.WriteHeader(); err != nil {
		t.Fatalf("WriteHeader error = %v", err)
	}
	if err := writer.WritePacket(time.Unix(1, 2000), []byte("abc")); err != nil {
		t.Fatalf("WritePacket error = %v", err)
	}
	if buf.Len() != 24+16+3 {
		t.Fatalf("buf.Len() = %d; want %d", buf.Len(), 43)
	}
}

func feedENHSequence(r *Reconstructor, symbols ...byte) [][]byte {
	var out [][]byte
	for _, symbol := range symbols {
		out = append(out, r.Feed(ENHFrame{Command: ENHResReceived, Data: symbol})...)
	}
	return out
}
