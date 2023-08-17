package middlewares

import (
	"roulette-api-server/config"
	"roulette-api-server/models/schema"
	"roulette-api-server/services"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func IsClientAuthenticated(c *gin.Context) {
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		var clientUser schema.OAuthClient
		if err := config.DB.
			Raw("select * from o_auth_clients where name = ? and secret = ?", user, password).
			Scan(&clientUser).
			Error; err != nil {
			services.Unauthorized(c, "Authentication Failed!", nil)
			return
		}
		c.Next()
	} else {
		services.Unauthorized(c, "Authentication Failed!", nil)
	}
}

func IsUserAuthenticated(c *gin.Context) {
	if c.Request.Header.Get("Authorization") == "" {
		services.Unauthorized(c, "Authentication Failed!", nil)
		return
	}

	req := c.Request.Header.Get("Authorization")
	var access schema.OAuthAccessToken
	if err := config.DB.
		Raw(
			"select * from o_auth_access_tokens where access_token = ? and expired_at > ? and revoked = ?",
			strings.Split(req, " ")[1],
			time.Now().Format("2006-01-02 15:04:05"),
			false,
		).
		Scan(&access).
		Error; err != nil {
		services.Unauthorized(c, "Invalid Token!", nil)
		return
	}
	c.Next()
}
