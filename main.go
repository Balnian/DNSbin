package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

const (
	baseDomaine = ".pts.com"
)

func main() {
	startDNS()
}

func startDNS() {
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}

type handler struct{}

func (srv *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	domain := msg.Question[0].Name
	stripdomain := strings.TrimRight(domain, baseDomaine)
	idx := strings.LastIndex(stripdomain, ".")
	uid := stripdomain[idx+1:]
	payload := stripdomain[:idx]
	fmt.Println(dns.TypeToString[r.Question[0].Qtype], ":", payload, uid)

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A:   net.ParseIP("127.0.0.1"),
		})
	case dns.TypeAAAA:
		msg.Authoritative = true
		msg.Answer = append(msg.Answer, &dns.AAAA{
			Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0},
			AAAA: net.ParseIP("::1"),
		})
	// Return empty data ()
	default:
		msg.Authoritative = true
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: 0},
		})
	}
	w.WriteMsg(&msg)
}
