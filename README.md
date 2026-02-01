# Go WhatsApp Proxy

A pure Go implementation of the WhatsApp Proxy, rewritten from the original Docker/HAProxy solution.

## Features
- **Pure Go**: No external dependencies like Docker or HAProxy required.
- **SOCKS5 Upstream**: Route all traffic through an external SOCKS5 proxy using the `SOCKS5_PROXY` env var.
- **Automatic TLS**: Generates self-signed certificates for port 443 on startup.
- **PROXY Protocol**: Support for sending PROXY v1 headers to WhatsApp backends.
- **Port Parity**: Implements all standard ports from the official WhatsApp proxy (80, 443, 5222, 587, 7777).

## Installation
```bash
go build -o whatsapp-proxy
```

## Configuration
Use environment variables to configure the proxy:
- `SOCKS5_PROXY`: (Optional) Upstream SOCKS5 address (e.g. `127.0.0.1:1080`).
- `LISTEN_ADDR`: (Optional) Address to bind to (default: `0.0.0.0`).

## Port Mappings
| Port | Backend Target | Description |
|------|----------------|-------------|
| 80 | g.whatsapp.net:80 | HTTP Proxy |
| 443 | g.whatsapp.net:5222 | HTTPS (TLS Terminated) |
| 5222 | g.whatsapp.net:5222 | XMPP |
| 587 | whatsapp.net:443 | Media Proxy |
| 7777 | whatsapp.net:443 | Media Proxy |
| 8199 | /stats | Stats & Healthcheck |

## License
MIT
