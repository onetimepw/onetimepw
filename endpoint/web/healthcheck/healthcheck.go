package healthcheck

import "github.com/onetimepw/onetimepw/build"

type HealthCheck struct {
	Uptime   int
	Version  string
	Release  string
	Strategy StrategyError
	checkers []Checker
}

type Checker func() CheckerResult

type CheckerResult struct {
	Checker        string `json:"checker"`
	Status         bool   `json:"status"`
	Message        string `json:"message,omitempty"`
	AdditionalInfo string `json:"additional_info,omitempty"`
}

func (h *HealthCheck) AddChecker(c Checker) {
	h.checkers = append(h.checkers, c)
}

func (h *HealthCheck) Check() (result map[string]interface{}, hasProblem bool) {
	// if have not strategy => add default strategy
	if h.Strategy == nil {
		h.Strategy = StrategyErrorOne{}
	}

	// run all checkers
	checks := make([]CheckerResult, len(h.checkers))

	for i, c := range h.checkers { // @todo make goroutines
		checks[i] = c()
	}

	hasProblem = h.Strategy.setError(checks)

	// prepare result
	result = map[string]interface{}{
		"release": build.Release,
		"version": build.Version,
		"uptime":  h.Uptime,
		"status":  "ok",
	}

	if hasProblem {
		result["status"] = "fail"
	}

	if len(checks) > 0 {
		result["results"] = checks
	}

	return result, hasProblem
}
