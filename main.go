package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/AYGA2K/todo_list/database"
	"github.com/AYGA2K/todo_list/routes"
)

func main() {
	db := database.ConnectDB()

	// Create a new table named "books"
	booksTable := "books"
	columns := map[string]string{
		"id":             "INTEGER PRIMARY KEY AUTOINCREMENT",
		"title":          "VARCHAR(255)",
		"author":         "VARCHAR(100)",
		"published_date": "DATE",
		// ... add more columns as needed ...
	}
	table, err := db.CreateTable(booksTable, columns)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(table.Name)

	r := gin.Default()
	r.GET("/oauth/callback", routes.OAuthCallbackHandler)
	r.GET("/google-drive", routes.GoogleDriveHandler)
	r.Run(":8080")
}
