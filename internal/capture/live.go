// SPDX-License-Identifier: GPL-2.0-or-later

package capture

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var errCaptureSinkClosed = errors.New("capture sink closed")

type LiveCaptureConfig struct {
	ProxyEndpoint string
	FIFOPath      string
	Stream        string
	ExportPCAPNG  string
}

func RunLiveCapture(cfg LiveCaptureConfig) error {
	if cfg.ProxyEndpoint == "" {
		return errors.New("missing proxy endpoint")
	}
	if cfg.FIFOPath == "" {
		return errors.New("missing required --fifo")
	}
	stream := StreamKind(cfg.Stream)
	switch stream {
	case StreamENSEvents, StreamEbusFrames, StreamBoth:
	default:
		return errors.New("unsupported --stream: expected ens-events, ebus-frames, or both")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	conn, err := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", cfg.ProxyEndpoint)
	if err != nil {
		return err
	}
	defer conn.Close()

	writer, err := os.OpenFile(cfg.FIFOPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer writer.Close()

	pcapWriter := NewPCAPWriter(writer)
	if err := pcapWriter.WriteHeader(); err != nil {
		return err
	}

	var mirror *os.File
	if cfg.ExportPCAPNG != "" {
		mirror, err = os.Create(cfg.ExportPCAPNG)
		if err != nil {
			return err
		}
		defer mirror.Close()
		if err := NewPCAPWriter(mirror).WriteHeader(); err != nil {
			return err
		}
	}

	parser := &ENHParser{}
	reconstructor := &Reconstructor{}
	mirrorWriter := NewPCAPWriter(mirror)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			return err
		}

		frame, err := parser.Parse(conn)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-ctx.Done():
					return nil
				default:
					continue
				}
			}
			return err
		}

		now := time.Now()
		if stream == StreamENSEvents || stream == StreamBoth {
			record := BuildENSEventPayload(frame)
			if err := writePacketSet(now, record, pcapWriter, mirror, mirrorWriter); err != nil {
				if errors.Is(err, errCaptureSinkClosed) {
					return nil
				}
				return err
			}
		}

		if stream == StreamEbusFrames || stream == StreamBoth {
			for _, record := range reconstructor.Feed(frame) {
				if err := writePacketSet(now, record, pcapWriter, mirror, mirrorWriter); err != nil {
					if errors.Is(err, errCaptureSinkClosed) {
						return nil
					}
					return err
				}
			}
		}
	}
}

func writePacketSet(ts time.Time, payload []byte, primary *PCAPWriter, mirrorFile *os.File, mirror *PCAPWriter) error {
	if err := primary.WritePacket(ts, payload); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return errCaptureSinkClosed
		}
		return err
	}
	if mirrorFile != nil {
		if err := mirror.WritePacket(ts, payload); err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return nil
			}
			return err
		}
	}
	return nil
}
