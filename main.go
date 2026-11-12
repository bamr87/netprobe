package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Probe reports whether a TCP connection to addr succeeds within timeout.
func Probe(addr string, timeout time.Duration) error {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	return c.Close()
}

func main() {
	if len(os.Args) < 2 {
		// A missing argument is a usage error: report it on stderr and exit 2
		// so scripts can tell it apart from a probe result (0 open, 1 closed).
		fmt.Fprintln(os.Stderr, "netprobe 0.1.0 — usage: netprobe host:port")
		os.Exit(2)
	}
	if err := Probe(os.Args[1], 2*time.Second); err != nil {
		fmt.Println("closed:", err)
		os.Exit(1)
	}
	fmt.Println("open:", os.Args[1])
}
