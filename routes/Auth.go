package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	googledrive "github.com/AYGA2K/todo_list/google_drive"
)

func OAuthCallbackHandler(c *gin.Context) {
	authCode := c.Query("code")
	fmt.Println(authCode)
	if authCode == "" {
		c.String(http.StatusBadRequest, "Authorization code not found in the request.")
		return
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to read client secret file: %v", err))
		return
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to parse client secret file to config: %v", err))
		return
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to retrieve token from web: %v", err))
		return
	}
	// Save the token to a file
	tokFile := "token.json"
	googledrive.SaveToken(tokFile, tok)

	c.String(http.StatusOK, "Authorization successful. Token saved.")
}

func GoogleDriveHandler(c *gin.Context) {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := googledrive.GetClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	f, err := os.Open("test.txt")
	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}

	defer f.Close()
	dir, err := googledrive.CreateDir(srv, "My Folder", "root")
	if err != nil {
		panic(fmt.Sprintf("Could not create dir: %v\n", err))
	}

	file, err := googledrive.CreateFile(srv, f.Name(), f.Name(), f, dir.Id)
	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	fmt.Printf("File '%s' successfully uploaded in '%s' directory", file.Name, dir.Name)
	c.IndentedJSON(http.StatusOK, "successful")
}
