package main

import (
	"clam/api"
	"clam/templates"
	"clam/templates/pages"
	"clam/templates/podcastHtml"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/gofiber/fiber/v2"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(c *fiber.Ctx) error {

	// Define template functions.
	metaTags := pages.MetaTags(
		"Podcast List",
		"Find the best podcasts",
	)
	bodyContent := pages.BodyContent(
		"Podcast List",
		"Find the best podcasts",
	)

	// Define template handler.
	templateHandler := templ.Handler(
		templates.Layout(
			"Welcome to example!", // define title text
			metaTags, bodyContent,
		),
	)

	// Render template layout.
	return adaptor.HTTPHandler(templateHandler)(c)

}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"status": "ok",
	})
}

func getPodcastsAPIHandler(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" || c.Get("HX-Request") != "true" {
		// If not, return HTTP 400 error.
		return fiber.NewError(fiber.StatusBadRequest, "non-htmx request")
	}

	search := c.Query("search")
	page := c.Query("page")

	// Get a list of podcasts.
	podcasts, err := api.GetPodcasts(search, page)

	if err != nil {
		// If there was an error, return HTTP 500 error.
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	// Return the list of podcasts as JSON.
	// return c.JSON(podcasts)
	templateHandler := templ.Handler(
		podcastHtml.PodcastList(podcasts),
	)
	return adaptor.HTTPHandler(templateHandler)(c)
}
