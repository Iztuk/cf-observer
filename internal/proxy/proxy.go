package proxy

import (
	"cf-observer/internal/config"
	"fmt"
	"log"
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
		host := host

		if host.Upstream == nil {
			return nil, fmt.Errorf("host %q has nil upstream", key)
		}

		rp := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(host.Upstream)
				pr.SetXForwarded()
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
