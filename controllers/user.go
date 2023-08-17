package controllers

import (
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/validations"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserFetchAll(c *gin.Context) {
	var user []schema.User
	err := models.UserFetchAll(&user)
	if err != nil {
		services.NotAcceptable(c, "Something went wrong!", err)
	} else {
		services.Success(c, nil, user)
	}
}

func UserFetchSingle(c *gin.Context) {
	userId := c.Param("id")

	var user schema.User
	err := models.UserFetchSingle(&user, userId)
	if err != nil {
		services.NotAcceptable(c, "Something went wrong!", err)
	} else {
		services.Success(c, nil, user)
	}
}

func UserCreate(c *gin.Context) {
	var request validations.UserCreate
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}

	request.Password = services.MD5Hash(request.Password)
	saveErr := models.UserCreate(&request)
	if saveErr != nil {
		services.NotAcceptable(c, "Something went wrong!", saveErr)
	} else {
		services.Success(c, "Created", request)
	}
}

func UserUpdate(c *gin.Context) {
	userId := c.Param("id")

	var request validations.UserUpdate
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}

	updateErr := models.UserUpdate(&request, userId)
	if updateErr != nil {
		services.NotAcceptable(c, "Something went wrong!", updateErr)
	} else {
		services.Success(c, "Updated", request)
	}
}

func UserDelete(c *gin.Context) {
	userId := c.Param("id")

	var user schema.User
	err := models.UserDelete(&user, userId)
	if err != nil {
		services.NotAcceptable(c, "Something went wrong!", err)
	} else {
		services.Success(c, "Deleted", nil)
	}
}
