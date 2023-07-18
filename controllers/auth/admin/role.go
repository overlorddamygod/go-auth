package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/overlorddamygod/go-auth/models"
	"github.com/overlorddamygod/go-auth/utils/response"
)

type CreateRoleParams struct {
	Name string `json:"name" binding:"required"`
}

func (a *AdminController) CreateRole(c *gin.Context) {
	var roleParams CreateRoleParams

	if err := c.Bind(&roleParams); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	// check if role exists
	var role models.Role
	result := a.db.Where("name = ?", roleParams.Name).First(&role)

	if result.Error == nil {
		response.BadRequest(c, "role already exists")
		return
	}

	result = a.db.Create(&models.Role{
		Name: roleParams.Name,
	})

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.Created(c, "role created")
}

type DeleteRoleParams struct {
	Id int `json:"id" binding:"required"`
}

func (a *AdminController) DeleteRole(c *gin.Context) {
	var deleteRoleParams DeleteRoleParams

	if err := c.Bind(&deleteRoleParams); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	result := a.db.Delete(&models.Role{}, "ud = ?", deleteRoleParams.Id)

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.Ok(c, "role deleted")
}

type AddRoleParams struct {
	UserId uuid.UUID `json:"user_id" binding:"required"`
	RoleId int       `json:"role_id" binding:"required"`
}

func (a *AdminController) AddRoleToUser(c *gin.Context) {
	var addRoleParams AddRoleParams

	if err := c.Bind(&addRoleParams); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	// add role to user_roles and check if user and role exists
	result := a.db.Create(&models.UserRole{
		UserID: addRoleParams.UserId,
		Type:   addRoleParams.RoleId,
	})

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.Created(c, "role added to user")
}

type RemoveRoleFromUserParams struct {
	UserId uuid.UUID `json:"user_id" binding:"required"`
	RoleId int       `json:"role_id" binding:"required"`
}

func (a *AdminController) RemoveRoleFromUser(c *gin.Context) {
	var RemoveRoleFromUserParams RemoveRoleFromUserParams

	if err := c.Bind(&RemoveRoleFromUserParams); err != nil {
		response.BadRequest(c, "invalid params")
		return
	}

	// add role to user_roles and check if user and role exists
	result := a.db.Delete(&models.UserRole{}, "user_id = ? AND role_id = ?", RemoveRoleFromUserParams.UserId, RemoveRoleFromUserParams.RoleId)

	if result.Error != nil {
		response.ServerError(c, "server error")
		return
	}

	response.Ok(c, "role removed from user")
}
