package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type Request struct {
	URL string `json:"url"`
}

// Enable CORS for requests from http://localhost:8081
func enableCors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin == "http://localhost:8081" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

// Handler function to generate thumbnail
func generateThumbnail(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx, fullScreenshot(req.URL, &buf))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Write the screenshot data to a temporary file
	tmpfile, err := ioutil.TempFile("./static", "screenshot-*.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the relative path of the temporary file within the static directory
	relativePath, err := filepath.Rel("./static", tmpfile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert relativePath to URL path using ToSlash for correct URL formatting
	thumbnailURL := "/static/" + filepath.ToSlash(relativePath)

	c.JSON(http.StatusOK, gin.H{
		"thumbnailUrl": thumbnailURL,
	})
}

// Task to capture a full screenshot using chromedp
func fullScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	r := gin.Default()

	// CORS handling middleware
	r.Use(func(c *gin.Context) {
		enableCors(c)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Serve static files using http.FileServer
	staticPath, _ := filepath.Abs("./static")
	r.GET("/static/*filepath", func(c *gin.Context) {
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))).ServeHTTP(c.Writer, c.Request)
	})

	// Endpoint to generate thumbnails
	r.POST("/generate", generateThumbnail)

	// Run the server on port 8080
	r.Run(":8080")
}

















