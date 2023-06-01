package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils/response"
	"gorm.io/gorm"
)

func (a *AdminController) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		response.BadRequest(c, "email is required")
		return
	}

	var user models.User

	result := a.db.First(&user, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.BadRequest(c, "user with that email not found")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	response.WithCustomStatusAndMessage(c, http.StatusOK, gin.H{
		"error": false,
		"user":  user,
	})
}

func (a *AdminController) GetUsersPaginated(c *gin.Context) {
	var users []models.User

	pageQ := c.Query("page")
	limitQ := c.Query("limit")
	page := 1
	limit := 10

	if pageQ != "" {
		p, err := strconv.Atoi(pageQ)

		if err == nil {
			if p > 0 {
				page = p
			}
		}
	}

	if limitQ != "" {
		l, err := strconv.Atoi(limitQ)

		if err == nil {
			if (l > 0) && (l <= 100) {
				limit = l
			}
		}
	}

	var count int64 = 0

	result := a.db.Limit(limit).Offset((page - 1) * limit).Find(&users).Count(&count)

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.WithCustomStatusAndMessage(c, http.StatusOK, gin.H{
		"error":     false,
		"users":     users,
		"page":      page,
		"limit":     limit,
		"totalPage": count / int64(limit),
	})
}

func (a *AdminController) GetAllTokens(c *gin.Context) {
	var users []models.RefreshToken

	result := a.db.Find(&users)

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.WithCustomStatusAndMessage(c, http.StatusOK, gin.H{
		"error": false,
		"users": users,
	})
}

func (a *AdminController) DeleteUser(c *gin.Context) {
	var user models.User

	result := a.db.First(&user, "id = ?", c.Param("id"))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.BadRequest(c, "user with that id not found")
			return
		}
		response.ServerError(c, "server error")
		return
	}

	result = a.db.Unscoped().Delete(&user)

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.Ok(c, "user deleted")
}
