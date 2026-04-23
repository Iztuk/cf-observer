package proxy

import (
	"cf-observer/internal/config"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func NewProxyManager(hosts map[string]config.Host, logger *log.Logger) (*ProxyManager, error) {
	pm := &ProxyManager{
		Hosts:  make(map[string]*ProxyTarget),
		Logger: logger,
	}

	for key, host := range hosts {
		if host.Upstream == nil {
			return nil, fmt.Errorf("host %q has nil upstream", key)
		}

		rp := httputil.NewSingleHostReverseProxy(host.Upstream)

		originalDirector := rp.Director
		rp.Director = func(r *http.Request) {
			originalHost := r.Host
			originalProto := "http"
			if r.TLS != nil {
				originalProto = "https"
			}

			originalDirector(r)

			r.Header.Set("X-Forwarded-Host", originalHost)
			r.Header.Set("X-Forwarded-Proto", originalProto)
		}

		rp.ModifyResponse = func(r *http.Response) error {
			return nil
		}

		rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "bad gateway", http.StatusBadGateway)
		}

		pm.Hosts[strings.ToLower(key)] = &ProxyTarget{
			Upstream: host.Upstream,
			Proxy:    rp,
			Logger:   logger,
		}
	}

	return pm, nil
}
