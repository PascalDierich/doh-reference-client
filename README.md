# DNS-over-HTTP reference Client

After I got interested in DNS-over-HTTP I looked at a few implementations.
Obviously most of them are pretty complex and do a lot of things, so I started to take a look at the IETF draft "DNS Queries over HTTPS (DoH)"[1]
and tried to implement it by my own. This is the result :).<p>

Each implementation detail is covered by a comment with the corresponding section in the IETF draft which (I hope) makes this project to a good reference.
Note that this is written after the 14th edition of named draft.<p>

While programming, I got more interested in the RFC 1035 [2] (DNS) as well, which lead to the result to only use code I fully understand. A consequence of it is that only the _A RDATA format_ of an _RR_ is read.<p>

Any kind of contribution to fix this or other issues is welcome.

### Usage
```
$ ./doh-reference-client -h
Usage of ./doh-reference-client:
  -address string
    	host address to resolve
  -method string
    	http method to use. Select "GET" or "POST" (default "GET")
  -server string
    	DoH server address (default "https://mozilla.cloudflare-dns.com/dns-query")
```

[1]: https://tools.ietf.org/html/draft-ietf-doh-dns-over-https-14 <br>
[2]: https://tools.ietf.org/html/rfc1035 <br>
