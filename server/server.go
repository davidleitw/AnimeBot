package server

import "github.com/gin-gonic/gin"

func AnimeBotServer() *gin.Engine {
	server := gin.Default()
	server.POST("/callback", callbackHandler)

	return server
}
