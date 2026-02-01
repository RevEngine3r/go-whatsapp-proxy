# Go WhatsApp Proxy

A pure Go implementation of the WhatsApp Proxy, rewritten from the original Docker/HAProxy solution.

## Features
- **Pure Go**: No external dependencies like Docker or HAProxy required (except for building).
- **SOCKS5 Upstream**: Optional support to route all traffic through an upstream SOCKS5 proxy.
- **Port Flexibility**: Supports all standard WhatsApp ports (80, 443, 5222, 587, 7777).
- **TLS Support**: Built-in TLS termination for HTTPS ports with automatic self-signed certificate generation.

## Usage
Run the binary with optional environment variables:
- `SOCKS5_PROXY`: Address of the upstream SOCKS5 proxy (e.g., `127.0.0.1:1080`).
- `LISTEN_ADDR`: Local address to bind to (default: `0.0.0.0`).
