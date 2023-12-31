package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	gowebly "github.com/gowebly/helpers"
)

// runServer runs a new HTTP server with the loaded environment variables.
func runServer() error {
	// Validate environment variables.
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7001"))
	if err != nil {
		return err
	}
	readTimeout, err := strconv.Atoi(gowebly.Getenv("BACKEND_READ_TIMEOUT", "5"))
	if err != nil {
		return err
	}
	writeTimeout, err := strconv.Atoi(gowebly.Getenv("BACKEND_WRITE_TIMEOUT", "10"))
	if err != nil {
		return err
	}

	// Create a new server instance with options from environment variables.
	// For more information, see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	config := fiber.Config{
		Views:        html.NewFileSystem(http.Dir("./templates"), ".html"),
		ViewsLayout:  "main",
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}

	// Create a new Fiber server.
	server := fiber.New(config)

	// Add Fiber middlewares.
	server.Use(logger.New())

	// Add rate limiting
	server.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// Handle static files.
	server.Static("/static", "./static")

	// Handle index page view.
	server.Get("/", indexViewHandler)

	// Handle health.
	server.Get("/health", healthHandler)

	// Handle API endpoints.
	server.Get("/api/podcasts", getPodcastsAPIHandler)

	return server.Listen(fmt.Sprintf(":%d", port))
}
