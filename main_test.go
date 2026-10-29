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
