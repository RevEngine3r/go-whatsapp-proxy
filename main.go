package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/proxy"
)

type Config struct {
	Socks5Upstream string
	ListenAddr     string
	StatsPort      int
}

type Backend struct {
	Target    string
	SendProxy bool
	IsTLS     bool
}

var backends = map[int]Backend{
	80:   {Target: "g.whatsapp.net:80", SendProxy: true},
	443:  {Target: "g.whatsapp.net:5222", SendProxy: true, IsTLS: true},
	5222: {Target: "g.whatsapp.net:5222", SendProxy: true},
	587:  {Target: "whatsapp.net:443", SendProxy: false},
	7777: {Target: "whatsapp.net:443", SendProxy: false},
}

func main() {
	cfg := Config{
		Socks5Upstream: os.Getenv("SOCKS5_PROXY"),
		ListenAddr:     os.Getenv("LISTEN_ADDR"),
		StatsPort:      8199,
	}
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = "0.0.0.0"
	}

	var dialer proxy.Dialer = proxy.Direct
	if cfg.Socks5Upstream != "" {
		d, err := proxy.SOCKS5("tcp", cfg.Socks5Upstream, nil, proxy.Direct)
		if err != nil {
			log.Fatalf("failed to create socks5 dialer: %v", err)
		}
		dialer = d
		log.Printf("Using upstream SOCKS5 proxy: %s", cfg.Socks5Upstream)
	}

	cert, err := GenerateSelfSignedCert()
	if err != nil {
		log.Fatalf("failed to generate cert: %v", err)
	}

	// Stats endpoint
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "WhatsApp Proxy (Go) is running.\n")
		})
		statsAddr := fmt.Sprintf("%s:%d", cfg.ListenAddr, cfg.StatsPort)
		log.Printf("Stats server listening on %s", statsAddr)
		if err := http.ListenAndServe(statsAddr, mux); err != nil {
			log.Printf("stats server error: %v", err)
		}
	}()

	var wg sync.WaitGroup
	for port, backend := range backends {
		wg.Add(1)
		go func(p int, b Backend) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%d", cfg.ListenAddr, p)
			var ln net.Listener
			var err error

			if b.IsTLS {
				tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
				ln, err = tls.Listen("tcp", addr, tlsCfg)
			} else {
				ln, err = net.Listen("tcp", addr)
			}

			if err != nil {
				log.Printf("failed to listen on %s: %v", addr, err)
				return
			}
			log.Printf("Listening on %s -> %s (ProxyProtocol: %v, TLS: %v)", addr, b.Target, b.SendProxy, b.IsTLS)

			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Printf("accept error on %d: %v", p, err)
					continue
				}
				go handleConnection(conn, b, dialer)
			}
		}(port, backend)
	}

	wg.Wait()
}

func handleConnection(clientConn net.Conn, b Backend, dialer proxy.Dialer) {
	defer clientConn.Close()

	targetConn, err := dialer.Dial("tcp", b.Target)
	if err != nil {
		log.Printf("failed to dial target %s: %v", b.Target, err)
		return
	}
	defer targetConn.Close()

	if b.SendProxy {
		// PROXY v1 header
		srcAddr, ok1 := clientConn.RemoteAddr().(*net.TCPAddr)
		dstAddr, ok2 := clientConn.LocalAddr().(*net.TCPAddr)
		if ok1 && ok2 {
			proxyHeader := fmt.Sprintf("PROXY TCP4 %s %s %d %d\r\n",
				srcAddr.IP.String(), dstAddr.IP.String(), srcAddr.Port, dstAddr.Port)
			targetConn.Write([]byte(proxyHeader))
		}
	}

	errChan := make(chan error, 2)
	copyFunc := func(dst, src net.Conn) {
		_, err := io.Copy(dst, src)
		errChan <- err
	}

	go copyFunc(targetConn, clientConn)
	go copyFunc(clientConn, targetConn)

	<-errChan
}
