package elblog

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	file, err := os.Open("data.log")
	if err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(file)
	logs := []Log{}
	for dec.More() {
		log, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		logs = append(logs, *log)
	}

	expected := []Log{
		// legacy log entries
		Log{
			Type: "http",
			Time: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
				return t
			}(),
			Name: "my-loadbalancer",
			From: &net.TCPAddr{
				IP:   net.ParseIP("192.168.131.39"),
				Port: 2817,
			},
			To: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.1"),
				Port: 80,
			},
			RequestProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("73µs")
				return d
			}(),
			BackendProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("1.048ms")
				return d
			}(),
			ResponseProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("57µs")
				return d
			}(),
			ELBStatusCode:     http.StatusOK,
			BackendStatusCode: "200",
			ReceivedBytes:     0,
			SentBytes:         29,
			Request:           "GET http://www.example.com:80/ HTTP/1.1",
			UserAgent:         "curl/7.38.0",
			SSLCipher:         "-",
			SSLProtocol:       "-",
		},
		Log{
			Type: "https",
			Time: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
				return t
			}(),
			Name: "my-loadbalancer",
			From: &net.TCPAddr{
				IP:   net.ParseIP("192.168.131.39"),
				Port: 2817,
			},
			To: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.1"),
				Port: 80,
			},
			RequestProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			BackendProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("2ms")
				return d
			}(),
			ResponseProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			ELBStatusCode:       http.StatusOK,
			BackendStatusCode:   "200",
			ReceivedBytes:       145,
			SentBytes:           1396,
			Request:             "GET https://www.example.com:443/ HTTP/1.1",
			UserAgent:           "-",
			SSLCipher:           "ECDHE-RSA-AES128-GCM-SHA256",
			SSLProtocol:         "TLSv1.2",
			TargetGroupARN:      "arn:aws:elasticloadbalancing:us-west-1:123456789012:targetgroup/aaaaa-aaaaa-ABCDEFGHIJKL/1111111111111111",
			TraceID:             "Root=1-11111111-111111111111111111111111",
			DomainName:          "www.example.com",
			ChosenCertARN:       "session-reused",
			MatchedRulePriority: "0",
		},
		// entries from june 2019
		Log{
			Type: "http",
			Time: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
				return t
			}(),
			Name: "my-loadbalancer",
			From: &net.TCPAddr{
				IP:   net.ParseIP("192.168.131.39"),
				Port: 2817,
			},
			To: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.1"),
				Port: 80,
			},
			RequestProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("73µs")
				return d
			}(),
			BackendProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("1.048ms")
				return d
			}(),
			ResponseProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("57µs")
				return d
			}(),
			ELBStatusCode:       http.StatusOK,
			BackendStatusCode:   "200",
			ReceivedBytes:       0,
			SentBytes:           29,
			Request:             "GET http://www.example.com:80/ HTTP/1.1",
			UserAgent:           "curl/7.38.0",
			SSLCipher:           "-",
			SSLProtocol:         "-",
			TargetGroupARN:      "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
			TraceID:             "Root=1-58337262-36d228ad5d99923122bbe354",
			DomainName:          "-",
			ChosenCertARN:       "-",
			MatchedRulePriority: "0",
			RequestCreationTime: "2018-07-02T22:22:48.364000Z",
			ActionsExecuted:     "forward",
			RedirectURL:         "-",
			ErrorReason:         "-",
		},
		Log{
			Type: "https",
			Time: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
				return t
			}(),
			Name: "my-loadbalancer",
			From: &net.TCPAddr{
				IP:   net.ParseIP("192.168.131.39"),
				Port: 2817,
			},
			To: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.1"),
				Port: 80,
			},
			RequestProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			BackendProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("2ms")
				return d
			}(),
			ResponseProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			ELBStatusCode:       http.StatusOK,
			BackendStatusCode:   "200",
			ReceivedBytes:       145,
			SentBytes:           1396,
			Request:             "GET https://www.example.com:443/ HTTP/1.1",
			UserAgent:           "-",
			SSLCipher:           "ECDHE-RSA-AES128-GCM-SHA256",
			SSLProtocol:         "TLSv1.2",
			TargetGroupARN:      "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
			TraceID:             "Root=1-58337281-1d84f3d73c47ec4e58577259",
			DomainName:          "www.example.com",
			ChosenCertARN:       "arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012",
			MatchedRulePriority: "1",
			RequestCreationTime: "2018-07-02T22:22:48.364000Z",
			ActionsExecuted:     "authenticate,forward",
			RedirectURL:         "-",
			ErrorReason:         "-",
		},
		// entries with other fields
		Log{
			Type: "https",
			Time: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
				return t
			}(),
			Name: "my-loadbalancer",
			From: &net.TCPAddr{
				IP:   net.ParseIP("192.168.131.39"),
				Port: 2817,
			},
			To: &net.TCPAddr{
				IP:   net.ParseIP("10.0.0.1"),
				Port: 80,
			},
			RequestProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			BackendProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("2ms")
				return d
			}(),
			ResponseProcessingTime: func() time.Duration {
				d, _ := time.ParseDuration("0s")
				return d
			}(),
			ELBStatusCode:       http.StatusOK,
			BackendStatusCode:   "200",
			ReceivedBytes:       145,
			SentBytes:           1396,
			Request:             "GET https://www.example.com:443/ HTTP/1.1",
			UserAgent:           "-",
			SSLCipher:           "ECDHE-RSA-AES128-GCM-SHA256",
			SSLProtocol:         "TLSv1.2",
			TargetGroupARN:      "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
			TraceID:             "Root=1-58337281-1d84f3d73c47ec4e58577259",
			DomainName:          "www.example.com",
			ChosenCertARN:       "arn:aws:acm:us-east-2:123456789012:certificate/12345678-1234-1234-1234-123456789012",
			MatchedRulePriority: "1",
			RequestCreationTime: "2018-07-02T22:22:48.364000Z",
			ActionsExecuted:     "authenticate,forward",
			RedirectURL:         "-",
			ErrorReason:         "-",
			OtherFields:         "future-entry-1 \"future-entry-2\" 3 future/entry/4",
		},
	}

	if !reflect.DeepEqual(expected, logs) {
		t.Fatalf("expected:\n	%v but got:\n	%v", expected, logs)
	}
}

func TestParse(t *testing.T) {
	cases := map[string]struct {
		given    string
		expected Log
	}{
		"basic": {
			given: `http 2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000073 0.001048 0.000057 200 200 0 29 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.38.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" "-" "-" 0 2018-07-02T22:22:48.364000Z "forward" "-" "-"`,
			expected: Log{
				Type: "http",
				Time: func() time.Time {
					t, _ := time.Parse(time.RFC3339, "2015-05-13T23:39:43.945958Z")
					return t
				}(),
				Name: "my-loadbalancer",
				From: &net.TCPAddr{
					IP:   net.ParseIP("192.168.131.39"),
					Port: 2817,
				},
				To: &net.TCPAddr{
					IP:   net.ParseIP("10.0.0.1"),
					Port: 80,
				},
				RequestProcessingTime: func() time.Duration {
					d, _ := time.ParseDuration("73µs")
					return d
				}(),
				BackendProcessingTime: func() time.Duration {
					d, _ := time.ParseDuration("1.048ms")
					return d
				}(),
				ResponseProcessingTime: func() time.Duration {
					d, _ := time.ParseDuration("57µs")
					return d
				}(),
				ELBStatusCode:       http.StatusOK,
				BackendStatusCode:   "200",
				ReceivedBytes:       0,
				SentBytes:           29,
				Request:             "GET http://www.example.com:80/ HTTP/1.1",
				UserAgent:           "curl/7.38.0",
				SSLCipher:           "-",
				SSLProtocol:         "-",
				TargetGroupARN:      "arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067",
				TraceID:             "Root=1-58337262-36d228ad5d99923122bbe354",
				DomainName:          "-",
				ChosenCertARN:       "-",
				MatchedRulePriority: "0",
				RequestCreationTime: "2018-07-02T22:22:48.364000Z",
				ActionsExecuted:     "forward",
				RedirectURL:         "-",
				ErrorReason:         "-",
			},
		},
	}

	for hint, c := range cases {
		t.Run(hint, func(t *testing.T) {
			got, err := Parse([]byte(c.given))
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}

			if !reflect.DeepEqual(*got, c.expected) {
				t.Errorf("expected:\n	%v but got:\n	%v", c.expected, *got)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	expected := 100
	buf := buffor(expected)
	dec := NewDecoder(buf)
	got := make([]*Log, 0, expected)
	for dec.More() {
		log, err := dec.Decode()
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		got = append(got, log)
	}
	if len(got) != expected {
		t.Errorf("wrong length, expected %d but got %d", expected, len(got))
	}
}

var benchLog Log

func BenchmarkParse(b *testing.B) {
	data := []byte(`http 2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000073 0.001048 0.000057 200 200 0 29 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.38.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" "-" "-" 0 2018-07-02T22:22:48.364000Z "forward" "-" "-"`)
	for n := 0; n < b.N; n++ {
		log, err := Parse(data)
		if err != nil {
			b.Fatalf("unexpected error: %s", err.Error())
		}
		benchLog = *log
	}
}

func buffor(max int) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	for i := 0; i < max; i++ {
		buf.WriteString(`http 2015-05-13T23:39:43.945958Z my-loadbalancer 192.168.131.39:2817 10.0.0.1:80 0.000073 0.001048 0.000057 200 200 0 29 "GET http://www.example.com:80/ HTTP/1.1" "curl/7.38.0" - - arn:aws:elasticloadbalancing:us-east-2:123456789012:targetgroup/my-targets/73e2d6bc24d8a067 "Root=1-58337262-36d228ad5d99923122bbe354" "-" "-" 0 2018-07-02T22:22:48.364000Z "forward" "-" "-"`)
		buf.WriteRune('\n')
	}
	return buf
}

func BenchmarkParse_NonParallel(b *testing.B) {
	buf := buffor(100000)
	b.ResetTimer()

	for n := 0; n <= b.N; n++ {
		b.StopTimer()
		buff := *buf
		scanner := bufio.NewScanner(&buff)
		scanner.Split(bufio.ScanLines)
		b.StartTimer()

		if scanner.Scan() {
			log, err := Parse(scanner.Bytes())
			if err != nil {
				b.Fatalf("unexpected error: %s", err.Error())
			}

			benchLog = *log
		}
	}
}

func BenchmarkParse_Parallel(b *testing.B) {
	buf := buffor(100000)
	parallelism := runtime.NumCPU() * 10
	b.ResetTimer()

	for n := 0; n <= b.N; n++ {
		b.StopTimer()

		buff := *buf
		scanner := bufio.NewScanner(&buff)
		scanner.Split(bufio.ScanLines)

		in := make(chan []byte)
		out := make(chan *Log)
		done := make(chan error, parallelism+1)

		for i := 0; i < parallelism; i++ {
			go func(in <-chan []byte, out chan<- *Log, done chan<- error) {
				for b := range in {
					log, err := Parse(b)
					if err != nil {
						done <- err
					}
					out <- log
				}
				done <- nil
			}(in, out, done)
		}

		go func(out <-chan *Log, done chan<- error) {
			for log := range out {
				benchLog = *log
			}
			done <- nil
		}(out, done)

		b.StartTimer()

		if scanner.Scan() {
			in <- scanner.Bytes()
		}
		close(in)

		kill := 0
	DoneLoop:
		for err := range done {
			if err != nil {
				b.Fatalf("unexpected error: %s", err.Error())
			}
			kill++
			if kill == parallelism {
				close(out)
				break DoneLoop
			}
		}
	}
}
