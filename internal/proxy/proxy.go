package proxy

import (
	"log"
	"net/http/httputil"
	"net/url"
)

type ProxyManager struct {
	Hosts  map[string]*ProxyTarget
	Logger *log.Logger
}

type ProxyTarget struct {
	Name     string
	Upstream *url.URL
	Proxy    *httputil.ReverseProxy

	Logger *log.Logger
}
