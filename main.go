package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

const timeout = 3 * time.Second

type test struct {
	name    string
	address string
}

func main() {
	tests := []test{
		{"IPv4", "1.1.1.1"},
		{"IPv6", "[2606:4700:4700::1111]"},
	}
	results := make([]chan string, len(tests))
	for i, t := range tests {
		c := make(chan string)
		results[i] = c
		go func() {
			c <- testDNS(t)
		}()
	}
	for i, c := range results {
		t := tests[i]
		fmt.Println(t.name, ":", <-c)
	}
}

func testDNS(t test) string {
	r := new(net.Resolver)
	r.PreferGo = true
	r.StrictErrors = true
	r.Dial = func(ctx context.Context, network, _ string) (net.Conn, error) {
		d := new(net.Dialer)
		d.FallbackDelay = -1
		return d.DialContext(ctx, "udp", t.address+":53")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := r.LookupHost(ctx, "cloudflare.com")
	if err != nil {
		return "FAIL: " + err.Error()
	} else {
		return "OK"
	}
}
