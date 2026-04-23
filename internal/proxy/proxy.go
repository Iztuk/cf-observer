package proxy

import (
	"cf-observer/internal/config"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ProxyManager struct {
	Hosts  map[string]*ProxyTarget
	Logger *log.Logger
}

type ProxyTarget struct {
	Upstream *url.URL
	Proxy    *httputil.ReverseProxy

	Logger *log.Logger
}

type Observation struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`

	RequestID string `json:"request_id"`

	Host     string `json:"host"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Query    string `json:"query"`
	Upstream string `json:"upstream"`

	Status     int   `json:"status"`
	DurationMs int64 `json:"duration_ms"`

	Error string `json:"error,omitempty"`

	RequestHeaders  map[string][]string `json:"request_headers,omitempty"`
	ResponseHeaders map[string][]string `json:"response_headers,omitempty"`
}

func NewProxyManager(hosts map[string]config.Host, logger *log.Logger) (*ProxyManager, error) {
	pm := &ProxyManager{
		Hosts:  make(map[string]*ProxyTarget),
		Logger: logger,
	}

	for key, host := range hosts {
		h := host

		if host.Upstream == nil {
			return nil, fmt.Errorf("host %q has nil upstream", key)
		}

		rp := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(h.Upstream)
				pr.SetXForwarded()

				requestID := getOrCreateRequestID(pr)

				obs := &Observation{
					Timestamp:      time.Now().UTC(),
					Event:          "request_started",
					RequestID:      requestID,
					Host:           pr.In.Host,
					Method:         pr.In.Method,
					Path:           pr.In.URL.Path,
					Query:          pr.In.URL.RawQuery,
					Upstream:       h.Upstream.String(),
					RequestHeaders: cloneHeader(pr.In.Header),
				}

				writeObservation(logger, obs)
			},
		}

		pm.Hosts[strings.ToLower(key)] = &ProxyTarget{
			Upstream: host.Upstream,
			Proxy:    rp,
			Logger:   logger,
		}
	}

	return pm, nil
}

func (pm *ProxyManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := normalizeHost(r.Host)

	target, ok := pm.Hosts[host]

	if !ok {
		pm.Logger.Printf("no route found for host=%s rawHost=%s", host, r.Host)
		http.NotFound(w, r)
		return
	}

	pm.Logger.Printf("routing host=%s to upstream=%s", host, target.Upstream.String())
	target.Proxy.ServeHTTP(w, r)
}

func normalizeHost(host string) string {
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err == nil {
			return strings.ToLower(h)
		}
	}
	return strings.ToLower(host)
}

func getOrCreateRequestID(r *httputil.ProxyRequest) string {
	if id := r.In.Header.Get("X-Request-ID"); id != "" {
		return id
	}

	if r.In.Header == nil {
		r.In.Header = make(http.Header)
	}

	var b [16]byte
	id := time.Now().UTC().Format("20060102150405.000000000")
	if _, err := rand.Read(b[:]); err == nil {
		id = hex.EncodeToString(b[:])
	}
	r.In.Header.Set("X-Request-ID", id)

	return id
}

func cloneHeader(h http.Header) map[string][]string {
	out := make(map[string][]string, len(h))
	for k, v := range h {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

func writeObservation(logger *log.Logger, obs *Observation) {
	b, err := json.Marshal(obs)
	if err != nil {
		logger.Printf(`{"message":"failed to marshal observation","error":%q}`, err.Error())
		return
	}
	logger.Print(string(b))
}
