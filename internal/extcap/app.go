// SPDX-License-Identifier: GPL-2.0-or-later

package extcap

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebus-extcap/internal/capture"
)

const (
	interfaceName = "helianthus-ebus-proxy"
	version       = "0.1.0"
)

type App struct {
	stdout io.Writer
}

func NewApp(stdout io.Writer) *App {
	return &App{stdout: stdout}
}

func (app *App) Run(args []string) error {
	fs := flag.NewFlagSet("helianthus-ebus-extcap", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		listInterfaces bool
		listDLTs       bool
		listConfig     bool
		extcapVersion  string
		captureMode    bool
		iface          string
		fifo           string
		proxyEndpoint  string
		transport      string
		stream         string
		exportPCAPNG   string
		srcFilter      string
		dstFilter      string
		opcodeFilter   string
		familyFilter   string
	)

	fs.BoolVar(&listInterfaces, "extcap-interfaces", false, "")
	fs.BoolVar(&listDLTs, "extcap-dlts", false, "")
	fs.BoolVar(&listConfig, "extcap-config", false, "")
	fs.StringVar(&extcapVersion, "extcap-version", "", "")
	fs.BoolVar(&captureMode, "capture", false, "")
	fs.StringVar(&iface, "extcap-interface", interfaceName, "")
	fs.StringVar(&fifo, "fifo", "", "")
	fs.StringVar(&proxyEndpoint, "proxy-endpoint", "", "")
	fs.StringVar(&transport, "transport", "ens", "")
	fs.StringVar(&stream, "stream", "both", "")
	fs.StringVar(&exportPCAPNG, "export-pcapng", "", "")
	fs.StringVar(&srcFilter, "src", "", "")
	fs.StringVar(&dstFilter, "dst", "", "")
	fs.StringVar(&opcodeFilter, "opcode", "", "")
	fs.StringVar(&familyFilter, "family", "", "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	switch {
	case listInterfaces:
		return app.listInterfaces()
	case listDLTs:
		return app.listDLTs(iface)
	case listConfig:
		return app.listConfig(iface)
	case captureMode:
		return app.capture(captureConfig{
			iface:         iface,
			fifo:          fifo,
			proxyEndpoint: proxyEndpoint,
			transport:     transport,
			stream:        stream,
			exportPCAPNG:  exportPCAPNG,
			srcFilter:     srcFilter,
			dstFilter:     dstFilter,
			opcodeFilter:  opcodeFilter,
			familyFilter:  familyFilter,
		})
	case extcapVersion != "":
		return app.printVersion()
	default:
		return app.listInterfaces()
	}
}

func (app *App) printVersion() error {
	_, err := fmt.Fprintf(app.stdout, "extcap {version=%s}{display=Helianthus eBUS passive capture}\n", version)
	return err
}

func (app *App) listInterfaces() error {
	_, err := fmt.Fprintf(app.stdout,
		"extcap {version=%s}{display=Helianthus eBUS passive capture}\n"+
			"interface {value=%s}{display=Helianthus eBUS via helianthus-ebusd-proxy}\n",
		version,
		interfaceName,
	)
	return err
}

func (app *App) listDLTs(iface string) error {
	if iface != interfaceName {
		return fmt.Errorf("unsupported extcap interface %q", iface)
	}
	_, err := fmt.Fprintln(app.stdout,
		"dlt {number=147}{name=helianthus-ebus}{display=USER0}",
	)
	return err
}

func (app *App) listConfig(iface string) error {
	if iface != interfaceName {
		return fmt.Errorf("unsupported extcap interface %q", iface)
	}
	lines := []string{
		"arg {number=0}{call=--proxy-endpoint}{display=Proxy endpoint}{type=string}{required=true}{tooltip=ENS northbound endpoint of helianthus-ebusd-proxy}",
		"arg {number=1}{call=--transport}{display=Transport}{type=selector}{default=ens}",
		"value {arg=1}{value=ens}{display=ENS}",
		"arg {number=2}{call=--stream}{display=Stream kind}{type=selector}{default=both}",
		"value {arg=2}{value=ens-events}{display=ENS events}",
		"value {arg=2}{value=ebus-frames}{display=eBUS frames}",
		"value {arg=2}{value=both}{display=Both streams}",
		"arg {number=3}{call=--export-pcapng}{display=Mirror to pcapng}{type=string}{required=false}",
		"arg {number=4}{call=--src}{display=Filter source address}{type=string}{required=false}",
		"arg {number=5}{call=--dst}{display=Filter target address}{type=string}{required=false}",
		"arg {number=6}{call=--opcode}{display=Filter opcode}{type=string}{required=false}",
		"arg {number=7}{call=--family}{display=Filter family}{type=string}{required=false}",
	}
	_, err := fmt.Fprintln(app.stdout, strings.Join(lines, "\n"))
	return err
}

type captureConfig struct {
	iface         string
	fifo          string
	proxyEndpoint string
	transport     string
	stream        string
	exportPCAPNG  string
	srcFilter     string
	dstFilter     string
	opcodeFilter  string
	familyFilter  string
}

func (app *App) capture(cfg captureConfig) error {
	if cfg.iface != interfaceName {
		return fmt.Errorf("unsupported extcap interface %q", cfg.iface)
	}
	if cfg.transport != "ens" {
		return fmt.Errorf("unsupported transport %q: v1 is proxy-only ENS", cfg.transport)
	}
	if cfg.proxyEndpoint == "" {
		return errors.New("missing required --proxy-endpoint")
	}
	if cfg.fifo == "" {
		return errors.New("missing required --fifo")
	}
	return capture.RunLiveCapture(capture.LiveCaptureConfig{
		ProxyEndpoint: cfg.proxyEndpoint,
		FIFOPath:      cfg.fifo,
		Stream:        cfg.stream,
		ExportPCAPNG:  cfg.exportPCAPNG,
	})
}
