package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

// closedAddr returns a loopback address with nothing listening on it: it
// reserves a free port, then releases it.
func closedAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestProbeOpen(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if err := Probe(l.Addr().String(), time.Second); err != nil {
		t.Fatalf("expected open: %v", err)
	}
}

func TestProbeClosed(t *testing.T) {
	addr := closedAddr(t)
	if err := Probe(addr, time.Second); err == nil {
		t.Fatalf("expected closed: probe of %s succeeded", addr)
	}
}

func TestProbeMissingPort(t *testing.T) {
	if err := Probe("127.0.0.1", time.Second); err == nil {
		t.Fatal("expected an error for an address without a port")
	}
}

// TestRun drives the CLI end to end: exit code, stdout and stderr per call.
// None of the usage-error cases dials, so the test needs no network.
func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	open, closed := l.Addr().String(), closedAddr(t)

	tests := []struct {
		name   string
		args   []string
		code   int
		stdout []string // substrings stdout must contain; nil means stdout must be empty
		stderr []string // substrings stderr must contain; nil means stderr must be empty
	}{
		{name: "open", args: []string{open}, code: 0, stdout: []string{"open: " + open + "\n"}},
		{name: "closed", args: []string{closed}, code: 1, stdout: []string{"closed: "}},
		{name: "no args", args: nil, code: 2, stderr: []string{"no address given", "usage: netprobe"}},
		{name: "missing port", args: []string{"example.com"}, code: 2, stderr: []string{"missing port", "usage: netprobe"}},
		{name: "empty port", args: []string{"example.com:"}, code: 2, stderr: []string{"missing port"}},
		{name: "empty host", args: []string{":443"}, code: 2, stderr: []string{"missing host"}},
		{name: "malformed address", args: []string{"::1"}, code: 2, stderr: []string{"too many colons"}},
		{name: "extra arg", args: []string{closed, "extra"}, code: 2, stderr: []string{`unexpected argument "extra"`, "usage: netprobe"}},
		{name: "unknown flag", args: []string{"--bogus", open}, code: 2, stderr: []string{"flag provided but not defined: -bogus", "usage: netprobe"}},
		{name: "help", args: []string{"--help"}, code: 0, stdout: []string{"usage: netprobe", "Exit codes:"}},
		{name: "short help", args: []string{"-h"}, code: 0, stdout: []string{"usage: netprobe", "Exit codes:"}},
		{name: "version", args: []string{"--version"}, code: 0, stdout: []string{"netprobe " + version + "\n"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != tt.code {
				t.Errorf("run(%q) = exit %d, want %d", tt.args, code, tt.code)
			}
			assertStream(t, "stdout", stdout.String(), tt.stdout)
			assertStream(t, "stderr", stderr.String(), tt.stderr)
		})
	}
}

// assertStream fails the test unless got contains every substring in want,
// or, when want is nil, unless got is empty.
func assertStream(t *testing.T, stream, got string, want []string) {
	t.Helper()
	if want == nil && got != "" {
		t.Errorf("%s = %q, want it empty", stream, got)
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("%s = %q, want it to contain %q", stream, got, w)
		}
	}
}
