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
		c.Error(err)
		return
	}

	book, err := h.service.CreateBook(c.Request.Context(), req.Title, req.Author, req.Price)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	logger.InfoWithRequestID("book created", requestID, map[string]interface{}{
		"book_id": book.ID,
	})

	c.Status(http.StatusCreated)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.GetByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	logger.InfoWithRequestID("book fetched", requestID, map[string]interface{}{
		"book_id": book.ID,
	})

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.service.GetAll(c)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	logger.InfoWithRequestID("all books fetched", requestID, map[string]interface{}{
		"total": len(books),
	})

	c.JSON(http.StatusOK, books)
}
