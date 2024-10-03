package web

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/onetimepw/onetimepw/domain"
	"github.com/onetimepw/onetimepw/usecase/api"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type Storage interface {
	Status() error
	Name() string
}

type Endpoint struct {
	fiber   *fiber.App
	config  domain.Config
	api     *api.API
	storage Storage
}

func New(
	config domain.Config,
	apiUC *api.API,
	storage Storage,
) (*Endpoint, error) {

	endpoint := &Endpoint{
		config:  config,
		api:     apiUC,
		storage: storage,
	}

	viewsPath := filepath.Join("./res", "public", "views")
	engine := html.New(viewsPath, ".html")

	if config.Env == "local" {
		engine.Reload(true)
	}

	app := fiber.New(fiber.Config{
		Views:        engine,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: endpoint.ErrorHandler,
	})

	app.Use(logger.New())

	limiterMiddleware := limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        30,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			// return c.Get("x-forwarded-for")
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Get("/health/check", endpoint.handlerCheck)

	apiGroup := app.Group("/api")
	{
		apiGroup.Use(limiterMiddleware).Post("/create", endpoint.handlerAPICreate)
		apiGroup.Post("view", endpoint.handlerAPIView)
	}
	app.Get("/view/:key", endpoint.handlerView)
	app.Post("/view/:key", endpoint.handlerView)

	endpoint.fiber = app

	return endpoint, nil
}

func (w *Endpoint) Run(ctx context.Context) error {
	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()

		zap.L().Info("Shutting down http server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return w.fiber.ShutdownWithContext(ctx)
	})

	g.Go(func() error {
		zap.L().Info("Starting http server",
			zap.Int("port", w.config.Port))
		port := strconv.Itoa(w.config.Port)

		if err := w.fiber.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	return g.Wait()
}
