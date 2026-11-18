package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// version is bumped by release-please; never edit it by hand.
const version = "0.2.0" // x-release-please-version

// probeTimeout bounds every TCP connect attempt.
const probeTimeout = 2 * time.Second

// Exit codes let a script tell a probe result apart from a mistake in the call.
const (
	exitOK     = 0 // the port is open, or --help / --version ran
	exitClosed = 1 // the TCP connection failed
	exitUsage  = 2 // the call itself was wrong; the usage is on stderr
)

// usage is the help text; %v is the probe timeout.
const usage = `usage: netprobe [-h | --help] [--version] host:port

Connect to host:port over TCP (timeout %v) and print the result on stdout:
"open: host:port" or "closed: <dial error>".

Examples:
  netprobe example.com:443
  netprobe '[::1]:22'    (IPv6 literals need brackets)

Exit codes:
  0  open, or --help / --version
  1  closed: the connection failed (refused, timed out, unresolvable host)
  2  usage error: no address, a missing host or port, a malformed address,
     an extra argument or an unknown flag
`

// Probe reports whether a TCP connection to addr succeeds within timeout.
func Probe(addr string, timeout time.Duration) error {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	return c.Close()
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the whole CLI: it parses args, writes the result to stdout or the
// error and usage to stderr, and returns the exit code, so tests can drive it
// without starting a process.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("netprobe", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // run reports parse errors itself, on the right stream
	showVersion := fs.Bool("version", false, "print the version and exit")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) { // -h, -help, --help
			printUsage(stdout)
			return exitOK
		}
		return usageError(stderr, err)
	}
	if *showVersion {
		fmt.Fprintf(stdout, "netprobe %s\n", version)
		return exitOK
	}
	switch fs.NArg() {
	case 0:
		return usageError(stderr, errors.New("no address given"))
	case 1:
	default:
		return usageError(stderr, fmt.Errorf("unexpected argument %q: give one host:port, with any flags before it", fs.Arg(1)))
	}
	addr := fs.Arg(0)
	if err := checkAddr(addr); err != nil {
		return usageError(stderr, err)
	}
	if err := Probe(addr, probeTimeout); err != nil {
		fmt.Fprintln(stdout, "closed:", err)
		return exitClosed
	}
	fmt.Fprintln(stdout, "open:", addr)
	return exitOK
}

// checkAddr reports why addr is not a host:port the probe can dial, so a
// malformed target is a usage error instead of a "closed" result.
func checkAddr(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		why := err.Error()
		var ae *net.AddrError
		if errors.As(err, &ae) {
			why = ae.Err // e.g. "missing port in address", without repeating addr
		}
		return fmt.Errorf("invalid address %q: %s", addr, why)
	}
	switch {
	case host == "":
		return fmt.Errorf("invalid address %q: missing host", addr)
	case port == "":
		return fmt.Errorf("invalid address %q: missing port", addr)
	}
	return nil
}

// usageError writes err and the usage to stderr and returns the usage exit code.
func usageError(stderr io.Writer, err error) int {
	fmt.Fprintf(stderr, "netprobe: %v\n\n", err)
	printUsage(stderr)
	return exitUsage
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usage, probeTimeout)
}
