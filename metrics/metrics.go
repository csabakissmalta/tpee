package metrics

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync/atomic"
	"time"
)

// metrics to be collected
// --- thiis is for tracking timeouts
// and other network-relatted issues ---
type HTTPMetrics struct {
	DNSLatency         time.Duration
	ConnectLatency     time.Duration
	TLSLatency         time.Duration
	WaitResponseHeader time.Duration
	TotalLatency       time.Duration
	ErrorCategory      string
	RetryCount         int
}

type ConnPoolStats struct {
	Idle   int
	Active int
}

var activeRequests int64

var transport *http.Transport

func classifyHTTPTimeout(err error) string {
	if err == nil {
		return ""
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "network_timeout"
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return "client_timeout"
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "no such host"):
		return "dns_failure"
	case strings.Contains(msg, "i/o timeout"):
		return "io_timeout"
	case strings.Contains(msg, "TLS handshake timeout"):
		return "tls_timeout"
	}

	return "other_error"
}

func InstrumentedRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		metrics := req.Context().Value("metrics").(*HTTPMetrics)
		start := time.Now()

		// track active and idle connections
		startRequest()
		defer endRequest()

		var dnsStart, connectStart, tlsStart, waitHeadersStart time.Time

		trace := &httptrace.ClientTrace{
			DNSStart: func(info httptrace.DNSStartInfo) {
				dnsStart = time.Now()
			},
			DNSDone: func(d httptrace.DNSDoneInfo) {
				metrics.DNSLatency = time.Since(dnsStart)
			},
			ConnectStart: func(network, addr string) {
				connectStart = time.Now()
			},
			ConnectDone: func(network, addr string, err error) {
				metrics.ConnectLatency = time.Since(connectStart)
			},
			TLSHandshakeStart: func() { tlsStart = time.Now() },
			TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
				metrics.TLSLatency = time.Since(tlsStart)
			},
			GotFirstResponseByte: func() {
				metrics.WaitResponseHeader = time.Since(waitHeadersStart)
			},
		}

		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		waitHeadersStart = time.Now()

		resp, err := rt.RoundTrip(req)

		metrics.TotalLatency = time.Since(start)
		metrics.ErrorCategory = classifyHTTPTimeout(err)

		return resp, err
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func startRequest() {
	atomic.AddInt64(&activeRequests, 1)
}

func endRequest() {
	atomic.AddInt64(&activeRequests, -1)
}

func GetTransport() *http.Transport {
	return transport
}

func GetConnPoolStats(tr *http.Transport) *ConnPoolStats {
	// Idle is approximated — exact values are not exposed
	active := atomic.LoadInt64(&activeRequests)
	idle := int64(tr.MaxIdleConnsPerHost) - active
	if idle < 0 {
		idle = 0
	}
	return &ConnPoolStats{Active: int(active), Idle: int(idle)}
}

func NewInstrumentedClient(redir bool) *http.Client {
	transport = &http.Transport{
		MaxIdleConns:          2000,
		MaxIdleConnsPerHost:   2000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 15 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
	}

	if redir {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: InstrumentedRoundTripper(transport),
		}
	} else {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: InstrumentedRoundTripper(transport),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

}
