package web

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (w *Endpoint) ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send custom error page
	switch code {
	case 404:
		return ctx.Render("404", fiber.Map{})
	default:
		errMsg := fmt.Sprintf("Error: %d", code)
		errMsg = errMsg + " " + http.StatusText(code)
		errMsg = errMsg + " " + err.Error()

		return ctx.Status(code).SendString(errMsg)
	}
}
