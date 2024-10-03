package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onetimepw/onetimepw/build"
	"github.com/onetimepw/onetimepw/endpoint/web/healthcheck"
	"time"
)

func (w *Endpoint) handlerCheck(c *fiber.Ctx) error {
	var checkers = []healthcheck.Checker{
		// состояние сервиса
		func() healthcheck.CheckerResult {
			return healthcheck.CheckerResult{
				Checker: "service",
				Status:  true,
			}
		},

		// состояние Storage
		func() healthcheck.CheckerResult {
			name := w.storage.Name()
			err := w.storage.Status()
			if err != nil {
				return healthcheck.CheckerResult{Checker: name, Status: false, Message: err.Error()}
			}

			return healthcheck.CheckerResult{
				Checker: name,
				Status:  true,
			}
		},
	}

	hCheck := &healthcheck.HealthCheck{
		Version:  build.Version,
		Release:  build.Release,
		Uptime:   int(time.Since(build.StartTime).Seconds()),
		Strategy: healthcheck.StrategyErrorOne{},
	}

	for _, checkerFunc := range checkers {
		hCheck.AddChecker(checkerFunc)
	}

	result, hasProblem := hCheck.Check()
	status := fiber.StatusOK
	if hasProblem {
		status = fiber.StatusInternalServerError
	}

	return c.Status(status).JSON(result)
}
