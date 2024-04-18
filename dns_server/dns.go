package dns_server

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type DnsServer struct {
	addr      string
	protocol  string
	proxyIP   string
	bypassUrl string
}

func NewDnsServer(addr string, proxy_ip string) *DnsServer {
	return &DnsServer{
		addr:      addr,
		proxyIP:   proxy_ip,
		protocol:  "udp",
		bypassUrl: "registry.docker.io.",
	}
}

func (s *DnsServer) Start() {
	// Define the DNS server address and port
	server := &dns.Server{Addr: s.addr, Net: s.protocol}

	// Define the DNS handler function
	server.Handler = dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)

		// Process each question in the DNS query
		for _, q := range r.Question {
			name := strings.ToLower(q.Name)

			// Check if the query is for registry.docker.io
			if strings.HasSuffix(name, s.bypassUrl) {
				// Redirect the query to your HTTP proxy IP
				proxyIP := s.proxyIP
				proxyIPs := net.ParseIP(proxyIP)
				if proxyIPs == nil {
					fmt.Printf("Invalid IP address for proxy: %s\n", proxyIP)
					continue
				}
				rr := &dns.A{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   proxyIPs,
				}
				m.Answer = append(m.Answer, rr)
			} else {
				// For other queries, return a default response
				// You can customize this behavior as needed
				m.SetRcode(r, dns.RcodeNameError)
			}
		}

		// Send the DNS response
		w.WriteMsg(m)
	})

	// Start the DNS server
	fmt.Printf("Starting DNS server on %s...", s.addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting DNS server: %s\n", err)
	}
}
