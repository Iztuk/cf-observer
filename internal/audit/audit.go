package audit

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"time"
)

type JobType string

const (
	RequestJobType  JobType = "request"
	ResponseJobType JobType = "response"
	FailureJobType  JobType = "failure"
)

type Metadata struct {
	RequestId  string
	Host       string
	Method     string
	Path       string
	Query      string
	Upstream   string
	Status     int
	Timestamp  time.Time
	DurationMs int64
}

type RequestJob struct {
	Type    JobType
	Meta    Metadata
	Headers http.Header
	Body    []byte
}

type ResponseJob struct {
	Type    JobType
	Meta    Metadata
	Headers http.Header
	Body    []byte
}

type FailureJob struct {
	Type  JobType
	Meta  Metadata
	Error string
}

func NewRequestJob(r *http.Request, upstream string, start time.Time) *RequestJob {
	requestId := getOrCreateRequestId(r)

	host := r.Header.Get("X-Original-Host")
	if host == "" {
		host = r.Host
	}

	return &RequestJob{
		Type: RequestJobType,
		Meta: Metadata{
			RequestId: requestId,
			Host:      host,
			Method:    r.Method,
			Path:      r.URL.Path,
			Query:     r.URL.RawQuery,
			Upstream:  upstream,
			Timestamp: start,
		},
		Headers: cloneHeader(r.Header),
	}
}

func NewResponseJob(r *http.Response, upstream string) *ResponseJob {
	requestId := getOrCreateRequestId(r.Request)

	start, _ := time.Parse(time.RFC3339Nano, r.Request.Header.Get("X-Request-Timestamp"))

	var duration int64
	if !start.IsZero() {
		duration = time.Since(start).Milliseconds()
	}

	host := r.Request.Header.Get("X-Original-Host")
	if host == "" {
		host = r.Request.Host
	}

	return &ResponseJob{
		Type: ResponseJobType,
		Meta: Metadata{
			RequestId:  requestId,
			Host:       host,
			Method:     r.Request.Method,
			Path:       r.Request.URL.Path,
			Query:      r.Request.URL.RawQuery,
			Upstream:   upstream,
			Status:     r.StatusCode,
			Timestamp:  time.Now().UTC(),
			DurationMs: duration,
		},
		Headers: cloneHeader(r.Header),
	}
}

func NewFailureJob(r *http.Request, upstream string, err error) *FailureJob {
	requestId := getOrCreateRequestId(r)

	status := http.StatusBadGateway
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		status = http.StatusGatewayTimeout
	}

	start, _ := time.Parse(time.RFC3339Nano, r.Header.Get("X-Request-Timestamp"))

	var duration int64
	if !start.IsZero() {
		duration = time.Since(start).Milliseconds()
	}

	host := r.Header.Get("X-Original-Host")
	if host == "" {
		host = r.Host
	}

	return &FailureJob{
		Type: FailureJobType,
		Meta: Metadata{
			RequestId:  requestId,
			Host:       host,
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.RawQuery,
			Upstream:   upstream,
			Status:     status,
			DurationMs: duration,
			Timestamp:  time.Now().UTC(),
		},
		Error: err.Error(),
	}
}

func getOrCreateRequestId(r *http.Request) string {
	if r.Header == nil {
		r.Header = make(http.Header)
	}

	id := r.Header.Get("X-Request-ID")
	if id == "" {
		id = newRequestId()
		r.Header.Set("X-Request-ID", id)
	}

	return id
}

func newRequestId() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}

	return time.Now().UTC().Format("20060102150405.000000000")
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
