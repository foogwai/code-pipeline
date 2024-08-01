package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/crseat/example-data-pipeline/internal/app"
	"github.com/crseat/example-data-pipeline/internal/domain"
)

// Handler represents an HTTP handler for managing requests.
type Handler struct {
	service *app.ProducerService
}

// NewHandler creates a new Handler with the provided ProducerService.
func NewHandler(service *app.ProducerService) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers the HTTP routes with the provided Echo instance.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/submit", h.handleSubmit)
}

// handleSubmit handles the HTTP POST request to the /submit endpoint.
// It binds and validates the JSON data from the request body, processes the post data,
// and returns the appropriate HTTP response.
func (h *Handler) handleSubmit(c echo.Context) error {
	var postData domain.PostData

	// Bind and validate the JSON data
	if err := c.Bind(&postData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	// Validate the struct
	if err := c.Validate(&postData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Process the post data
	if err := h.service.ProcessPostData(postData); err != nil {
		// In a production setting this would mean kafka was down which should really never happen if we configure it
		// correctly. Even so we would probably have some logic here to send the failed messages to an S3 bucket, so
		// we can re-drive them into kafka after outage is over.
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// No content response if data is valid
	return c.NoContent(http.StatusNoContent)
}
