package main

import (
	"net"
	"testing"
	"time"
)

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
	// Reserve a free port, then release it so nothing is listening there.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	if err := Probe(addr, time.Second); err == nil {
		t.Fatalf("expected closed: probe of %s succeeded", addr)
	}
}

func TestProbeMissingPort(t *testing.T) {
	if err := Probe("127.0.0.1", time.Second); err == nil {
		t.Fatal("expected an error for an address without a port")
	}
}
