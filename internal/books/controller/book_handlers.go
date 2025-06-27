package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/books/application"
	"github.com/mingtmt/book-store/internal/middleware"
	"github.com/mingtmt/book-store/pkg/logger"
)

type BookHandler struct {
	service *application.BookService
}

func NewBookHandler(service *application.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Price  string `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// Debug: log the raw request body for troubleshooting
		body, _ := c.GetRawData()
		logger.Error("failed to bind JSON", err, map[string]interface{}{
			"raw_body": string(body),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, make sure price is a string (e.g., \"39.99\")"})
		return
	}

	book, err := h.service.CreateBook(c.Request.Context(), req.Title, req.Author, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	logger.InfoWithRequestID("book created", requestID, map[string]interface{}{
		"book_id": book.ID,
	})

	c.JSON(http.StatusCreated, book)
}
