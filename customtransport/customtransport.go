package customtransport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type CustomTransport struct {
	Rtp       http.RoundTripper
	dialer    *net.Dialer
	connStart time.Time
	connEnd   time.Time
	reqStart  time.Time
	reqEnd    time.Time
}

type Option func(*CustomTransport)

func WithTimeout(tout time.Duration) Option {
	return func(tr *CustomTransport) {
		tr.dialer.Timeout = tout
	}
}

func WithKeepAliveDuration(ka time.Duration) Option {
	return func(tr *CustomTransport) {
		tr.dialer.KeepAlive = ka
	}
}

func WithDisabledSSLVerify(verify bool) Option {
	return func(tr *CustomTransport) {
		tr.Rtp = &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			Dial:                tr.dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}
	}
}

func NewTransport(option ...Option) *CustomTransport {
	tr := &CustomTransport{
		dialer: &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	}
	tr.Rtp = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                tr.dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	for _, o := range option {
		o(tr)
	}
	return tr
}

func (tr *CustomTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr.reqStart = time.Now()
	resp, err := tr.Rtp.RoundTrip(r)
	tr.reqEnd = time.Now()
	return resp, err
}

func (tr *CustomTransport) dial(network, addr string) (net.Conn, error) {
	tr.connStart = time.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = time.Now()
	return cn, err
}

func (tr *CustomTransport) ReqDuration() time.Duration {
	return tr.Duration() - tr.ConnDuration()
}

func (tr *CustomTransport) ConnDuration() time.Duration {
	return tr.connEnd.Sub(tr.connStart)
}

func (tr *CustomTransport) Duration() time.Duration {
	return tr.reqEnd.Sub(tr.reqStart)
}
