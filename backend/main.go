package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"time"
)

type Request struct {
	URL string `json:"url"`
}

func enableCors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin == "http://localhost:8081" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func generateThumbnail(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Enable CORS
	enableCors(c)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	// Define a custom wait time (2 seconds in milliseconds)
	waitTime := 2000

	err := chromedp.Run(ctx, fullScreenshot(req.URL, waitTime, &buf))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	tmpfile, err := ioutil.TempFile("", "screenshot-*.png")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(buf); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"thumbnailUrl": tmpfile.Name(),
	})
}

func fullScreenshot(urlstr string, waitTime int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(time.Duration(waitTime) * time.Millisecond), // Use the custom wait time here
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	r := gin.Default()

	// CORS handling middleware
	r.Use(func(c *gin.Context) {
		enableCors(c)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	r.POST("/generate", generateThumbnail)
	r.Run(":8080")
}









