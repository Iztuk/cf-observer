package audit

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ResponseStatusCodeRule struct{}

func (r ResponseStatusCodeRule) ID() RuleID {
	return RuleResponseStatusCodeNotDefined
}

func (r ResponseStatusCodeRule) Title() string {
	return "Response status code not defined"
}

func (r ResponseStatusCodeRule) AppliesTo() []JobType {
	return []JobType{ResponseJobType}
}

func (r ResponseStatusCodeRule) Check(ctx RuleContext, job Job, jobID string) ([]Finding, error) {
	responseJob, ok := job.(*ResponseJob)
	if !ok {
		return nil, nil
	}

	op, found := ctx.Contracts.FindMethod(
		responseJob.Meta.Host,
		responseJob.Meta.Method,
		responseJob.Meta.Path,
	)
	if !found {
		return nil, nil
	}

	status := strconv.Itoa(responseJob.Meta.Status)
	_, found = op.Responses[status]
	if found {
		return nil, nil
	}

	return []Finding{
		{
			ID:     uuid.NewString(),
			JobID:  jobID,
			RuleID: string(r.ID()),
			Title:  r.Title(),
			Message: fmt.Sprintf(
				"Response status code %d for %s %s is not defined in the API contract.",
				responseJob.Meta.Status,
				responseJob.Meta.Method,
				responseJob.Meta.Path,
			),
			CreatedAt: time.Now().UTC(),
		},
	}, nil
}
