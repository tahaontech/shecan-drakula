package main

import (
	"github.com/tahaontech/shecan-drakula/dns_server"
)

func main() {
	dnsServer := dns_server.NewDnsServer(":53", "<your_proxy_ip>")

	dnsServer.Start()
}
