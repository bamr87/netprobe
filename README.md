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

```bash
go vet ./... && go test ./...
```

## How it grows

The roadmap lives in [PRD.md](PRD.md) and [features/features.yml](features/features.yml) (ids `NETP-NNN`, most `status: planned`). Work arrives as pull requests from the hub's evolution and fan-out loops, the `@claude` handler, or the local `evolve` lane (off until `EVOLVE_ENABLED=true`). A human merges every PR. Rules: [docs/constitution.md](docs/constitution.md) · agent guide: [AGENTS.md](AGENTS.md).

## License

[MIT](LICENSE)
