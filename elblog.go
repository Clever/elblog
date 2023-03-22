package elblog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"unicode/utf8"
)

// Log ...
type Log struct {
	Type                   string
	Time                   time.Time
	Name                   string
	From, To               *net.TCPAddr
	RequestProcessingTime  time.Duration
	BackendProcessingTime  time.Duration
	ResponseProcessingTime time.Duration
	ELBStatusCode          int
	BackendStatusCode      string
	ReceivedBytes          int64
	SentBytes              int64
	Request                string
	UserAgent              string
	SSLCipher              string
	SSLProtocol            string
	TargetGroupARN         string
	TraceID                string
	DomainName             string
	ChosenCertARN          string
	MatchedRulePriority    string
	RequestCreationTime    string
	ActionsExecuted        string
	RedirectURL            string
	ErrorReason            string
	TargetPortList         string
	TargetStatusCodeList   string
	Classification         string
	ClassificationReason   string
	OtherFields            string
}

const numTokens = 30

// Parse ...
func Parse(b []byte) (log *Log, err error) {
	var (
		adv, i int
		code   int64
		dur    float64
		ip     int64
		tok    []byte
		parts  [][]byte
	)

	data := b[adv:]
	log = &Log{}
	i = 0
	for i < numTokens && adv < len(data) {
		data = data[adv:]
		adv, tok, err = scan(data)
		if err != nil {
			return nil, fmt.Errorf("unable to scan next token: %v", err)
		}
		switch i {
		case 0:
			log.Type = string(tok)
		case 1:
			log.Time, err = time.Parse(time.RFC3339Nano, string(tok))
		case 2:
			log.Name = string(tok)
		case 3:
			parts = bytes.Split(tok, []byte(":"))
			switch len(parts) {
			case 1:
				log.From = &net.TCPAddr{
					IP: net.ParseIP(string(parts[0])),
				}
			case 2:
				ip, err = strconv.ParseInt(string(parts[1]), 10, 32)
				log.From = &net.TCPAddr{
					IP:   net.ParseIP(string(parts[0])),
					Port: int(ip),
				}
			}
		case 4:
			parts = bytes.Split(tok, []byte(":"))
			switch len(parts) {
			case 1:
				log.To = &net.TCPAddr{
					IP: net.ParseIP(string(parts[0])),
				}
			case 2:
				ip, err = strconv.ParseInt(string(parts[1]), 10, 32)
				log.To = &net.TCPAddr{
					IP:   net.ParseIP(string(parts[0])),
					Port: int(ip),
				}
			}
		case 5:
			dur, err = strconv.ParseFloat(string(tok), 64)
			log.RequestProcessingTime = time.Duration(dur * 1000 * 1000 * 1000)
		case 6:
			dur, err = strconv.ParseFloat(string(tok), 64)
			log.BackendProcessingTime = time.Duration(dur * 1000 * 1000 * 1000)
		case 7:
			dur, err = strconv.ParseFloat(string(tok), 64)
			log.ResponseProcessingTime = time.Duration(dur * 1000 * 1000 * 1000)
		case 8:
			code, err = strconv.ParseInt(string(tok), 10, 32)
			log.ELBStatusCode = int(code)
		case 9:
			log.BackendStatusCode = string(tok)
		case 10:
			log.ReceivedBytes, err = strconv.ParseInt(string(tok), 10, 32)
		case 11:
			log.SentBytes, err = strconv.ParseInt(string(tok), 10, 32)
		case 12:
			log.Request = string(tok)
		case 13:
			log.UserAgent = string(tok)
		case 14:
			log.SSLCipher = string(tok)
		case 15:
			log.SSLProtocol = string(tok)
		case 16:
			log.TargetGroupARN = string(tok)
		case 17:
			log.TraceID = string(tok)
		case 18:
			log.DomainName = string(tok)
		case 19:
			log.ChosenCertARN = string(tok)
		case 20:
			log.MatchedRulePriority = string(tok)
		case 21:
			log.RequestCreationTime = string(tok)
		case 22:
			log.ActionsExecuted = string(tok)
		case 23:
			log.RedirectURL = string(tok)
		case 24:
			log.ErrorReason = string(tok)
		case 25:
			log.TargetPortList = string(tok)
		case 26:
			log.TargetStatusCodeList = string(tok)
		case 27:
			log.Classification = string(tok)
		case 28:
			log.ClassificationReason = string(tok)
		case 29:
			// we've scanned one token but we want to put everything remaining into OtherFields
			// (including the spaces and quotes)
			log.OtherFields = string(data)
			adv = len(data)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid field %q at index %d: %v", tok, i, err)
		}
		i++
	}
	return
}

// scan works like bufio.ScanWord (most of the code is taken from there),
// but treat everything between quotation marks also as a word.
func scan(data []byte) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	open := false
	trim := false
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != ' ' {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '"' {
			trim = true
			open = !open

		}
		if r == ' ' && !open {
			if trim {
				return i + width, data[start+1 : i-1], nil
			}
			return i + width, data[start:i], nil
		}
	}
	// We have a final, non-empty, non-terminated word. Return it.
	if len(data) > start {
		if trim && !open {
			return len(data), data[start+1 : len(data)-1], nil
		}
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

// Decoder ...
type Decoder struct {
	s     *bufio.Scanner
	token []byte
}

// NewDecoder allocates new Decoder object for given input.
func NewDecoder(r io.Reader) *Decoder {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	return &Decoder{
		s: s,
	}
}

// Decode scans input and parse into Log. It can return EOF if underlying scanner Scan method returns false.
func (d *Decoder) Decode() (*Log, error) {
	if d.token != nil {
		log, err := Parse(d.token)
		d.token = nil
		if err != nil {
			return nil, err
		}
		return log, nil
	}
	ok := d.s.Scan()
	if !ok {
		return nil, io.EOF
	}
	return Parse(d.s.Bytes())
}

// More return true if token is not empty or underlying scanner Scan method will return true.
func (d *Decoder) More() bool {
	if d.token != nil {
		return true
	}

	ok := d.s.Scan()
	if ok {
		d.token = d.s.Bytes()
	}
	return ok
}
