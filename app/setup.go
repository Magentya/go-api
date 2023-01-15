package app

import (
	"go-api/configs"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
)

func Setup() error {

	configs.Setup()

	// Server
	app := fiber.New(fiber.Config{
		// Custom settings
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber - Go",
		AppName:       "Backend",
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		BodyLimit:     4 * 1024 * 1024,
		Concurrency:   256 * 1024,
	})

	// Middleware

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "X-CSRF-Token",
		CookieSameSite: "Strict",
		CookieSecure:   true,
		Expiration:     3600,
		KeyGenerator:   utils.UUIDv4,
	}))

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	app.Get("/metrics/", monitor.New(monitor.Config{Title: "Metrics Page"}))

	// Routes

	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(
			fiber.Map{
				"message": "Hello, World!",
			},
		)
	})

	app.Listen(":" + os.Getenv("PORT"))

	return nil
}
