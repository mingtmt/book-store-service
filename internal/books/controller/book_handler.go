package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/books/application"
	"github.com/mingtmt/book-store/internal/middleware"
	"github.com/mingtmt/book-store/pkg/logger"
)

// BookRequest represents the payload for creating or updating a book
type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  string `json:"price"`
}

type BookHandler struct {
	service *application.BookService
}

func NewBookHandler(service *application.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with title, author, and price
// @Tags books
// @Accept json
// @Produce json
// @Param book body controller.BookRequest true "Book details"
// @Success 201 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /books [post]
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
	userID := c.GetString("userID")
	logger.InfoWithRequestID("book created", requestID, map[string]interface{}{
		"book_id": book.ID,
		"user_id": userID,
	})

	c.Status(http.StatusCreated)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Get a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.GetBookByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	userID := c.GetString("userID")
	logger.InfoWithRequestID("book fetched", requestID, map[string]interface{}{
		"book_id": book.ID,
		"user_id": userID,
	})

	c.JSON(http.StatusOK, book)
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get a list of all books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} domain.Book
// @Failure 500 {object} response.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.service.GetAllBooks(c)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	userID := c.GetString("userID")
	logger.InfoWithRequestID("all books fetched", requestID, map[string]interface{}{
		"total":   len(books),
		"user_id": userID,
	})

	c.JSON(http.StatusOK, books)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update a book's title, author, and price by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body controller.BookRequest true "Updated book details"
// @Success 200 {object} domain.Book
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Price  string `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	book, err := h.service.UpdateBook(c.Request.Context(), id, req.Title, req.Author, req.Price)
	if err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	userID := c.GetString("userID")
	logger.InfoWithRequestID("book updated", requestID, map[string]interface{}{
		"book_id": book.ID,
		"user_id": userID,
	})

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteBookByID(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	requestID := c.GetString(middleware.RequestIDKey)
	userID := c.GetString("userID")
	logger.InfoWithRequestID("book deleted", requestID, map[string]interface{}{
		"book_id": id,
		"user_id": userID,
	})

	c.Status(http.StatusNoContent)
}
