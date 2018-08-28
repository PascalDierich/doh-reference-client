package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pascaldierich/doh-reference-client/lib"
)

// media type for DoH is defined as following. (see: section 6)
const mediaType = "application/dns-message"

// "The DoH client SHOULD include an HTTP 'Accept' request header
// field to indicate what type of content can be understood in response."
// (see: section 4.1)
var setAcceptHeader = func(req *http.Request) {
	req.Header.Set("Accept", mediaType)
}

// Sends a DNS-over-HTTP POST request and returns the answered RR's
// or error if any.
func sendPOSTRequest(server, addr string) ([]*lib.RR, error) {
	payload, err := createDNSMessage(addr)
	if err != nil {
		return nil, err
	}

	// "When using the POST method the DNS query is included as the message
	// body of the HTTP request and the Content-Type request header field
	// indicates the media type of the message."
	// (see: secttion 4.1)
	req, err := http.NewRequest(http.MethodPost, server, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mediaType)
	setAcceptHeader(req)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// "A successful HTTP response with a 2xx status code is used for any
	// valid DNS response, regardless of the DNS response code."
	// (see: section 4.2.1)
	if resp.StatusCode != http.StatusOK {
		tmp := fmt.Sprintf("bad response Code: %v", resp.StatusCode)
		return nil, errors.New(tmp)
	}

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseMsg := &lib.Message{}
	err = lib.UnmarshalMessage(msg, responseMsg)
	if err != nil {
		return nil, err
	}

	return responseMsg.Answers, nil
}

// "When the HTTP method is GET the single variable 'dns' is defined as
// the content of the DNS request". (see: section 4.1)
const DNSRequestContent = "dns"

// Sends a DNS-over-HTTP GET request and returns the answered RR's
// or error if any.
func sendGETRequest(server, addr string) ([]*lib.RR, error) {
	payload, err := createDNSMessage(addr)
	if err != nil {
		return nil, err
	}

	// "When using the GET method, the data payload [...] MUST be encoded
	// with base64url...". (see: section 6)
	enc := base64.RawURLEncoding.EncodeToString(payload)

	// "...and then provided as a variable named 'dns' to the URI Template
	// expansion. (see: section 6)
	url := fmt.Sprintf("%s?%s=%s", server, DNSRequestContent, enc)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mediaType)
	setAcceptHeader(req)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// "A successful HTTP response with a 2xx status code is used for any
	// valid DNS response, regardless of the DNS response code."
	// (see: section 4.2.1)
	if resp.StatusCode != http.StatusOK {
		tmp := fmt.Sprintf("bad response Code: %v", resp.StatusCode)
		return nil, errors.New(tmp)
	}

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseMsg := &lib.Message{}
	err = lib.UnmarshalMessage(msg, responseMsg)
	if err != nil {
		return nil, err
	}

	return responseMsg.Answers, nil
}

// Returns raw DNS message.
func createDNSMessage(addr string) (raw []byte, err error) {
	queryMsg := &lib.Message{
		Header: lib.Header{
			// "In order to maximize cache friendliness, DoH clients using media
			// formats that include DNS ID, such as application/dns-message, SHOULD
			// use a DNS ID of 0 in every DNS request."
			// (see: section 4.1)
			ID:      0,
			QR:      0,
			Opcode:  lib.OpcodeQuery,
			QDCOUNT: 1,
			RD:      1,
		},
		Questions: []*lib.Question{
			{
				QNAME:  addr,
				QTYPE:  lib.QTypeA,
				QCLASS: lib.QClassIN,
			},
		},
	}

	raw, err = queryMsg.Marshal()
	return
}
