package audit

import (
	"time"

	"github.com/google/uuid"
)

type RuleID string

const (
	// RuleProxyUpstreamFailure applies when the proxy successfully matches a host
	// and route target, but cannot complete the upstream request due to a
	// non-timeout proxy error.
	//
	// Example:
	//   - upstream service is not running
	//   - connection is refused
	//   - upstream connection is reset
	//
	// This rule is evaluated against FailureJob values.
	RuleProxyUpstreamFailure RuleID = "proxy.upstream_failure"

	// RuleProxyUpstreamTimeout applies when the proxy successfully matches a host
	// and route target, but the upstream request exceeds the configured timeout.
	//
	// Example:
	//   - upstream does not accept the connection before Dial timeout
	//   - upstream accepts the connection but does not return response headers
	//     before ResponseHeaderTimeout
	//
	// This rule is evaluated against FailureJob values.
	RuleProxyUpstreamTimeout RuleID = "proxy.upstream_timeout"

	// RuleRequestPathDoesNotExist applies when an incoming request path cannot be
	// matched to any path defined in the OpenAPI contract for the selected host.
	//
	// Example:
	//   - request:  GET /admin
	//   - contract: /users, /users/{id}, /health
	//
	// This rule should run before method, content type, and body validation because
	// those checks require a matching contract path.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestPathDoesNotExist RuleID = "request.path_does_not_exist"

	// RuleRequestMethodNotAllowed applies when the request path exists in the
	// OpenAPI contract, but the specific HTTP method is not defined for that path.
	//
	// Example:
	//   - request:  DELETE /users
	//   - contract: GET /users and POST /users only
	//
	// This rule should run after path matching succeeds, but before request body
	// validation because the operation definition is needed for deeper checks.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestMethodNotAllowed RuleID = "request.method_not_allowed"

	// RuleRequestContentTypeNotAllowed applies when the request has a body and the
	// Content-Type header does not match any media type allowed by the OpenAPI
	// operation's requestBody.content map.
	//
	// Example:
	//   - request Content-Type: text/plain
	//   - contract allows: application/json
	//
	// This validates the declared media type only. It does not prove the body is
	// actually valid JSON, XML, multipart data, etc.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestContentTypeNotAllowed RuleID = "request.content_type_not_allowed"

	// RuleRequestBodyMissing applies when the OpenAPI operation declares a required
	// request body, but the captured request body is empty.
	//
	// Example:
	//   - contract: requestBody.required = true
	//   - request: POST /users with no body
	//
	// This rule should run after the path and method have been resolved to an
	// OpenAPI operation.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestBodyMissing RuleID = "request.body_missing"

	// RuleRequestBodyNotAllowed applies when the OpenAPI operation does not define
	// a requestBody, but the incoming request includes a body.
	//
	// Example:
	//   - contract: GET /health has no requestBody
	//   - request: GET /health with a JSON body
	//
	// This catches clients sending payloads to operations that are expected to be
	// bodyless.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestBodyNotAllowed RuleID = "request.body_not_allowed"

	// RuleRequestInvalidJSON applies when the OpenAPI operation expects a JSON
	// request body and the captured request body is not valid JSON.
	//
	// Example:
	//   - request Content-Type: application/json
	//   - body: {"email":
	//
	// This rule should only run when the resolved operation allows a JSON media
	// type such as application/json or application/*+json.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestInvalidJSON RuleID = "request.invalid_json"

	// RuleRequestInvalidBodyFormat applies when the request body does not match the
	// expected non-JSON media format declared by the OpenAPI contract.
	//
	// Example future uses:
	//   - multipart/form-data body cannot be parsed as multipart data
	//   - application/xml body cannot be parsed as XML
	//   - text/csv body cannot be parsed as CSV
	//
	// This is a generic extension point for media-type-specific validators beyond
	// JSON. It should run only after content type validation determines which media
	// type applies.
	//
	// This rule is evaluated against RequestJob values.
	RuleRequestInvalidBodyFormat RuleID = "request.invalid_body_format"
)

type Rule interface {
	ID() RuleID
	Title() string
	AppliesTo() []JobType
	Check(job Job, jobID string) ([]Finding, error)
}

type RuleEngine struct {
	rules []Rule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: getRules(),
	}
}

func getRules() []Rule {
	return []Rule{
		UpstreamFailureRule{},
		UpstreamTimeoutRule{},
	}
}

func (e *RuleEngine) Evaluate(job Job, jobID string) ([]Finding, error) {
	var findings []Finding

	for _, rule := range e.rules {
		if !ruleApplies(rule, job.JobType()) {
			continue
		}

		ruleFindings, err := rule.Check(job, jobID)
		if err != nil {
			return nil, err
		}

		findings = append(findings, ruleFindings...)
	}

	return findings, nil
}

func ruleApplies(rule Rule, jobType JobType) bool {
	for _, supported := range rule.AppliesTo() {
		if supported == jobType {
			return true
		}
	}

	return false
}

type UpstreamFailureRule struct{}

func (r UpstreamFailureRule) ID() RuleID {
	return RuleProxyUpstreamFailure
}

func (r UpstreamFailureRule) Title() string {
	return "Upstream request failed"
}

func (r UpstreamFailureRule) AppliesTo() []JobType {
	return []JobType{FailureJobType}
}

func (r UpstreamFailureRule) Check(job Job, jobID string) ([]Finding, error) {
	failureJob, ok := job.(*FailureJob)
	if !ok {
		return nil, nil
	}

	if failureJob.Meta.Status == 504 {
		return nil, nil
	}

	return []Finding{
		{
			ID:        uuid.NewString(),
			JobID:     jobID,
			RuleID:    string(r.ID()),
			Title:     r.Title(),
			Message:   failureJob.Error,
			CreatedAt: time.Now().UTC(),
		}}, nil
}

type UpstreamTimeoutRule struct{}

func (r UpstreamTimeoutRule) ID() RuleID {
	return RuleProxyUpstreamTimeout
}

func (r UpstreamTimeoutRule) Title() string {
	return "Upstream request timed out"
}

func (r UpstreamTimeoutRule) AppliesTo() []JobType {
	return []JobType{FailureJobType}
}

func (r UpstreamTimeoutRule) Check(job Job, jobID string) ([]Finding, error) {
	failureJob, ok := job.(*FailureJob)
	if !ok {
		return nil, nil
	}

	if failureJob.Meta.Status != 504 {
		return nil, nil
	}

	return []Finding{
		{
			ID:        uuid.NewString(),
			JobID:     jobID,
			RuleID:    string(r.ID()),
			Title:     r.Title(),
			Message:   "The upstream service did not respond before the configured timeout.",
			CreatedAt: time.Now().UTC(),
		},
	}, nil
}
