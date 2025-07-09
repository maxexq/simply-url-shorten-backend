package shotenhandler

import (
	"fmt"

	service "url-shortener/internal/service/shorten-service"

	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(s service.URLService) *URLHandler {
	return &URLHandler{s}
}

func (h *URLHandler) ShortenURL(c *fiber.Ctx) error {
	url := c.FormValue("url")
	if url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "URL is required"})
	}

	urlModel, err := h.service.Shorten(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to shorten URL"})
	}

	return c.JSON(fiber.Map{"short_url": fmt.Sprintf("http://localhost:8080/%s", urlModel.ShortCode)})
}

func (h *URLHandler) RedirectURL(c *fiber.Ctx) error {
	code := c.Params("code")
	originalURL, err := h.service.Resolve(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
	}
	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}
