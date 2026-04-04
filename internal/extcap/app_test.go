package extcap

import (
	"bytes"
	"strings"
	"testing"
)

func TestListInterfaces(t *testing.T) {
	var out bytes.Buffer
	app := NewApp(&out)
	if err := app.Run([]string{"--extcap-interfaces"}); err != nil {
		t.Fatalf("Run(--extcap-interfaces) error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "interface {value=helianthus-ebus-proxy}") {
		t.Fatalf("interfaces output missing expected interface:\n%s", got)
	}
}

func TestVersion(t *testing.T) {
	var out bytes.Buffer
	app := NewApp(&out)
	if err := app.Run([]string{"--extcap-version", "4.4"}); err != nil {
		t.Fatalf("Run(--extcap-version) error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "extcap {version=0.1.0}") {
		t.Fatalf("version output missing expected extcap version stanza:\n%s", got)
	}
}

func TestListConfig(t *testing.T) {
	var out bytes.Buffer
	app := NewApp(&out)
	if err := app.Run([]string{"--extcap-interface", "helianthus-ebus-proxy", "--extcap-config"}); err != nil {
		t.Fatalf("Run(--extcap-config) error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "--proxy-endpoint") {
		t.Fatalf("config output missing --proxy-endpoint:\n%s", got)
	}
	if !strings.Contains(got, "--stream") {
		t.Fatalf("config output missing --stream:\n%s", got)
	}
}

func TestCaptureRequiresFIFO(t *testing.T) {
	var out bytes.Buffer
	app := NewApp(&out)
	err := app.Run([]string{
		"--capture",
		"--extcap-interface", "helianthus-ebus-proxy",
		"--proxy-endpoint", "127.0.0.1:19001",
	})
	if err == nil {
		t.Fatal("Run(--capture) error = nil; want fifo error")
	}
	if !strings.Contains(err.Error(), "missing required --fifo") {
		t.Fatalf("capture error = %q", err)
	}
}
