# netprobe

> A Go network diagnostics toolkit with a terminal UI: ping, port scan, DNS, traceroute and HTTP probes.

One binary for the network questions you ask every week, with a live TUI dashboard built on Bubble Tea.

**Status:** foundation (phase 1). Part of the [bamr87 fleet](https://github.com/bamr87/bamr87).

## Goals

- Fast concurrent TCP port scanner and ICMP/TCP ping
- DNS lookup and traceroute views
- Bubble Tea TUI tying probes into one dashboard

## Skills learned

- Go concurrency (goroutines, contexts)
- Raw sockets and network protocols
- TUI design with Bubble Tea
- Cross-platform release builds

## Fleet tie-in

Diagnoses fleet hosts and services (fleet-infra cluster, fleet-observe endpoints) and can feed probe results to the activity log.

## Quick start

Requires Go 1.22 or newer (see `go.mod`).

```bash
go vet ./... && go test ./...   # the CI gate
go build -o netprobe .          # build the binary
```

## Usage

Today the binary does one thing: a TCP connect probe (NETP-001) against a single `host:port`, with a fixed 2-second timeout. Port scan, DNS, traceroute, HTTP probes and the TUI are planned (see [PRD.md](PRD.md)).

```bash
./netprobe example.com:443   # prints "open: example.com:443", exit 0
./netprobe 127.0.0.1:9       # prints "closed: <dial error>", exit 1
./netprobe example.com       # prints the error (missing port) and usage to stderr, exit 2
./netprobe                   # prints the error (no address given) and usage to stderr, exit 2
./netprobe --help            # prints usage and the exit codes to stdout, exit 0 (also -h)
./netprobe --version         # prints "netprobe <version>", exit 0
```

Flags go before the address. IPv6 literals need brackets: `./netprobe '[::1]:22'`. "closed" covers every dial failure, including a refused connection, a timeout and a name that does not resolve; the error text says which. A target that is not a `host:port` is never reported as "closed": it is a usage error, so a script can tell a mistake in the call from a closed port.

| Exit code | Meaning |
| --- | --- |
| 0 | TCP connection succeeded, or `--help` / `--version` ran |
| 1 | Connection failed (refused, timed out, unresolvable host); prints `closed: <dial error>` on stdout |
| 2 | Usage error (no address, a missing host or port, a malformed address, an extra argument or an unknown flag); prints the error and usage on stderr, nothing on stdout |

## How it grows

The roadmap lives in [PRD.md](PRD.md) and [features/features.yml](features/features.yml) (ids `NETP-NNN`, most `status: planned`). Work arrives as pull requests from the hub's evolution and fan-out loops, the `@claude` handler, or the local `evolve` lane (off until `EVOLVE_ENABLED=true`). A human merges every PR. Rules: [docs/constitution.md](docs/constitution.md) · agent guide: [AGENTS.md](AGENTS.md) · documentation index: [docs/README.md](docs/README.md).

## License

[MIT](LICENSE)
