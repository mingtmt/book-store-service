package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/common/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if any errors were added
		err := c.Errors.Last()
		if err != nil {
			// Check for AppError type
			if appErr, ok := err.Err.(*errors.AppError); ok {
				c.JSON(appErr.Status, gin.H{
					"error": gin.H{
						"code":    appErr.Code,
						"message": appErr.Message,
					},
				})
				return
			}

			// Fallback: unexpected error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "unexpected error occurred",
				},
			})
		}
	}
}
