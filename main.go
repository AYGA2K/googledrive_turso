package main

import (
	"github.com/gin-gonic/gin"

	"github.com/AYGA2K/todo_list/routes"
)

func main() {
	r := gin.Default()
	r.GET("/oauth/callback", routes.OAuthCallbackHandler)
	r.GET("/google-drive", routes.GoogleDriveHandler)
	r.Run(":8080")
}
