package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func response(c *gin.Context, err bool, status int, altmessage string, message string) {
	if message == "" {
		message = altmessage
	}
	c.JSON(status, gin.H{
		"error":   err,
		"message": message,
	})
}

func WithCustomStatusAndMessage(c *gin.Context, status int, message interface{}) {
	c.JSON(status, message)
}

func Ok(c *gin.Context, message string) {
	response(c, false, http.StatusOK, "Success", message)
}

func Created(c *gin.Context, message string) {
	response(c, false, http.StatusCreated, "created", message)
}

func ServerError(c *gin.Context, message string) {
	response(c, true, http.StatusInternalServerError, "internal server error", message)
}

func Unauthorized(c *gin.Context, message string) {
	response(c, true, http.StatusUnauthorized, "unauthorized", message)
}

func BadRequest(c *gin.Context, message string) {
	response(c, true, http.StatusBadRequest, "bad request", message)
}

func NotFound(c *gin.Context, message string) {
	response(c, true, http.StatusNotFound, "not found", message)
}
