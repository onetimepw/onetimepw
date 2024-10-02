package web

import (
	"app/build"
	"app/endpoint/web/healthcheck"
	"github.com/gofiber/fiber/v2"
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

		// состояние Redis
		func() healthcheck.CheckerResult {
			name := "redis"
			conn := w.redisClient.Conn()
			defer conn.Close()

			// Test the connection
			statusCmd := conn.Ping(c.Context())
			if statusCmd == nil {
				return healthcheck.CheckerResult{Checker: name, Status: false, Message: "can't run ping command on redis"}
			}

			err := statusCmd.Err()
			if err != nil {
				return healthcheck.CheckerResult{Checker: name, Status: false, Message: err.Error()}
			}

			_, err = statusCmd.Result()
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
