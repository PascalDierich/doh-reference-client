package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pascaldierich/doh-reference-client/lib"
)

var server = flag.String("server", "https://mozilla.cloudflare-dns.com/dns-query", "DoH server address")
var method = flag.String("method", "GET", "http method to use. Select \"GET\" or \"POST\"")
var address = flag.String("address", "", "host address to resolve")

func main() {
	flag.Parse()

	if *address == "" {
		log.Fatal("Set address to resolve")
	}

	switch *method {
	case "GET", "get":
		rrs, err := sendGETRequest(*server, *address)
		if err != nil {
			log.Fatal(err)
		}
		printIP(rrs)

	case "POST", "post":
		rrs, err := sendGETRequest(*server, *address)
		if err != nil {
			log.Fatal(err)
		}
		printIP(rrs)

	default:
		log.Fatal("Select http method. Use \"GET\" or \"POST\"")
	}
}

func printIP(rrs []*lib.RR) {
	// NOTE: At the moment we can only unmarshal A RDATA formats
	// as described in RFC 1035. Thats why we print out the last
	// RDATA section which most likely is an A RDATA format,
	// as the previous would be a CNAME.
	var ip string
	for _, rr := range rrs {
		ip = UnmarshalRDATA(rr.RDATA)
	}
	fmt.Printf("-> %v\n", ip)
}

// Unmarshals RDATA format as A RDATA format as described in RFC 1035.
func UnmarshalRDATA(data []byte) (ip string) {
	for _, b := range data {
		ip += fmt.Sprintf("%d.", int(b))
	}
	ip = strings.TrimSuffix(ip, ".")
	return
}
