package dnslogger

import (
	"log"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

//DNSQueryData contains Data from a DNS Requests
type DNSQueryData struct {
	DomainQuery string
	QueryType   string
	Addr        string
}

//DNSDataHandler is the function that handle the data
type DNSDataHandler func(dnsdata DNSQueryData)

var dnsLogFunc DNSDataHandler

//StartDNSLogger start a basic DNS server that will answer "localhost" to A and AAAA query
func StartDNSLogger(logfunc DNSDataHandler) {
	dnsLogFunc = logfunc
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

	dnsLogFunc(DNSQueryData{DomainQuery: msg.Question[0].Name, QueryType: dns.TypeToString[r.Question[0].Qtype], Addr: w.RemoteAddr().String()})
	domain := msg.Question[0].Name
	/*stripdomain := strings.TrimRight(domain, baseDomaine)
	idx := strings.LastIndex(stripdomain, ".")
	uid := stripdomain[idx+1:]
	payload := stripdomain[:idx]
	fmt.Println(dns.TypeToString[r.Question[0].Qtype], ":", payload, uid)*/

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
