package metrics

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
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

// AI-generated -- doesn't work

// func GetConnPoolStats(transport *http.Transport) ConnPoolStats {
// 	idle := 0
// 	active := 0

// 	transport.Range(func(key, value interface{}) bool {
// 		if pc, ok := value.(*http.Transport).IdleConnTimeout; ok {
// 			fmt.Println(pc)
// 		}
// 		return true
// 	})

// 	// For Go <1.22: need reflection (Go doesn't expose idle/active counts)
// 	// For now, you can measure:
// 	//   - total inflight requests (you track)
// 	//   - idle connections via Transport.IdleConnTimeout + MaxIdleConnsPerHost

// 	return ConnPoolStats{Idle: idle, Active: active}
// }

func NewInstrumentedClient(redir bool) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:          2000,
		MaxIdleConnsPerHost:   2000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
	}

	if redir {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: InstrumentedRoundTripper(tr),
		}
	} else {
		return &http.Client{
			Timeout:   30 * time.Second,
			Transport: InstrumentedRoundTripper(tr),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

}
