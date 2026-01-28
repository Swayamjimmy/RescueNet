# RescueNet (Comms / P2P Chat)

This repository contains a small Go project that demonstrates **local-network peer discovery** and **chat messaging** using **libp2p**:

- **Peer discovery**: mDNS (find peers on the same LAN)
- **Messaging**: GossipSub PubSub (broadcast messages to a “room” topic)
- **Optional HTTP API**: send a message and fetch messages via REST endpoints

The runnable program is currently `cmd/comms/main.go`, and the reusable P2P helpers live in `internal/p2p`.

## What this does

Each node:

- creates a libp2p host listening on a TCP port you provide (`--port`)
- discovers other peers using mDNS (`--same_string` is the rendezvous string that must match)
- joins a PubSub topic derived from the room name (`--room`)
- publishes messages typed into stdin (terminal)
- receives messages from peers, prints them to the terminal, and appends them to `cmd/comms/logs.txt`
- (optional) runs an HTTP server (defaults to `:3001`) that can publish messages and return the collected message list

## Project layout

- `cmd/comms/main.go`: CLI app (host + mDNS + pubsub + stdin loop + optional HTTP)
- `internal/p2p/host.go`: host creation (`CreateHost`)
- `internal/p2p/mdns.go`: mDNS discovery (`InitMDNS`)
- `internal/p2p/pubsub.go`: chat room topic/subscription (`JoinChatRoom`, `Publish`)

## Requirements

- Go 1.20+ (Go modules)
- A local network where mDNS works (same Wi‑Fi/LAN)

Notes:

- On **WSL2**, mDNS discovery can be unreliable depending on your Windows/WSL networking setup. If peers don’t discover each other, see [Troubleshooting](#troubleshooting).

## Install / build

From the repo root:

```bash
go mod download
go build ./...
```

## Run (terminal chat)

You need **at least 2 terminals**, usually on **two machines** on the same LAN, or two different WSL/host sessions where mDNS works.

### Terminal 1

```bash
cd cmd/comms
go run . --port 9000 --same_string xyz --room myroom --nick Swayam
```

### Terminal 2

```bash
cd cmd/comms
go run . --port 9001 --same_string xyz --room myroom --nick Anshika
```

Now type a line and press Enter in either terminal to publish it to the room.

## Run with HTTP API (optional)

Start a node with HTTP enabled:

```bash
cd cmd/comms
go run . --port 9000 --same_string xyz --room myroom --nick Swayam --enable-http true
```

### Endpoints

- `POST /send`
  - Body JSON: `{"message":"hello"}`
- `GET /messages`
  - Returns JSON array of strings collected in memory on that node

### Examples (curl)

Send a message:

```bash
curl -X POST http://localhost:3001/send \
  -H 'Content-Type: application/json' \
  -d '{"message":"hello from http"}'
```

Fetch messages stored in memory:

```bash
curl http://localhost:3001/messages
```

## CLI flags

- `--port` (required): TCP port for the libp2p host to listen on
- `--nick` (optional): nickname displayed with messages
- `--room` (optional): room name (maps to a PubSub topic)
- `--same_string` (recommended): mDNS rendezvous string; must match across peers to discover each other
- `--enable-http` (optional): start the HTTP server on `:3001`

## Files produced

When you run `cmd/comms`, it appends incoming messages to:

- `cmd/comms/logs.txt`

## Troubleshooting

### “error creating the host”

Common causes:

- **Missing `--port`**: it must be provided
- **Port already in use**: choose a different `--port` (e.g. 9001, 9002, …)

### Peers are not discovering each other

- Ensure:
  - both nodes use the **same** `--same_string`
  - both nodes are on the **same LAN**
  - UDP/mDNS is not blocked by firewall/router rules
- On WSL2:
  - discovery may fail depending on Windows firewall + WSL virtual NIC behavior
  - easiest workaround is to run one node on Windows native Go and one in WSL, or run both on Linux/macOS hosts on the same LAN

### HTTP requests fail in a browser

If you’re calling the HTTP server from a frontend, CORS headers must be correct. If you see browser CORS errors, check the header names returned by the server.

## Development notes

- PubSub topic naming currently prefixes the room name (see `internal/p2p/pubsub.go`).
- Messages are broadcast to everyone in the topic; there is no authentication/encryption beyond what libp2p provides by default.

## License

No license file is currently included. Add one (MIT/Apache-2.0/etc.) if you plan to share this publicly.