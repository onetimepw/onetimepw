package web

import (
	"github.com/gofiber/fiber/v2"
	"slices"
	"time"
)

type CreateRequest struct {
	Text     string `json:"text"`
	Password string `json:"password"`
	Duration string `json:"duration"`
}

func (w *Endpoint) handlerAPICreate(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	allowedDurations := []string{
		"15m",
		"30m",
		"1h",
		"2h",
		"3h",
	}

	// if duration is not allowed, set default "1h"
	if !slices.Contains(allowedDurations, req.Duration) {
		req.Duration = "1h"
	}

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		duration = time.Hour
	}

	key, pass, err := w.api.Create(req.Text, req.Password, duration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":       key,
		"password": pass,
	})
}

type ViewRequest struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

func (w *Endpoint) handlerAPIView(c *fiber.Ctx) error {
	var req ViewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	text, err := w.api.Get(req.Key, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"text": text,
	})
}

func (w *Endpoint) handlerView(c *fiber.Ctx) error {
	key := c.Params("key")

	if !w.api.Has(key) {
		return c.Status(fiber.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("view", fiber.Map{})
}
