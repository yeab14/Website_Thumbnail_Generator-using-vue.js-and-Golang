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

func generateThumbnail(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx, fullScreenshot(req.URL, 90, &buf))
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

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	r := gin.Default()
	r.POST("/generate", generateThumbnail)
	r.Run(":8080")
}

