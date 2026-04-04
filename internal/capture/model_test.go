package capture

import "testing"

func TestSupportedOpcodeCatalog(t *testing.T) {
	if got := len(SupportedSemanticEbusOpcodes); got != 46 {
		t.Fatalf("len(SupportedSemanticEbusOpcodes) = %d; want 46", got)
	}
	if got := len(SupportedSemanticVaillantOpcodes); got != 11 {
		t.Fatalf("len(SupportedSemanticVaillantOpcodes) = %d; want 11", got)
	}
	if !IsSupportedSemanticOpcode(0xB524) {
		t.Fatal("0xB524 should be supported")
	}
	if IsSupportedSemanticOpcode(0xB599) {
		t.Fatal("0xB599 should not be supported")
	}
}

func TestSemanticLabelFallback(t *testing.T) {
	if got := SemanticLabel(0xB509); got != "Vaillant B509 boiler registers" {
		t.Fatalf("SemanticLabel(0xB509) = %q", got)
	}
	if got := SemanticLabel(0x1234); got != "Raw opcode 0x1234" {
		t.Fatalf("SemanticLabel(0x1234) = %q", got)
	}
}
