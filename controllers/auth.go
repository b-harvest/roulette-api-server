package controllers

import (
	"roulette-api-server/models"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"roulette-api-server/validations"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthSignin(c *gin.Context) {
	var request validations.AuthSignin
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}

	var emailUser schema.User
	emailUserErr := models.UserFetchWithEmail(&emailUser, request.Email)
	if emailUserErr != nil {
		services.NotAcceptable(c, "You do not have an account with this email!", nil)
		return
	}

	if emailUser.Password != services.MD5Hash(request.Password) {
		services.NotAcceptable(c, "Invalid email or password!", nil)
		return
	}

	clientName, _, _ := c.Request.BasicAuth()
	tokens := services.GenerateTokens(strconv.Itoa(int(emailUser.Id)),  clientName)

	services.Success(c, "Welcome!", gin.H{"userInfo": emailUser, "tokens": tokens})
}

func AuthSignout(c *gin.Context) {
	res := services.Signout(strings.Split(c.Request.Header.Get("Authorization"), " ")[1])
	if res != true {
		services.BadRequest(c, "Something went wrong!", nil)
		return
	}
	services.Deleted(c, "Logout Successful!", nil)
	return
}
