// +build go1.8
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

const (
	TIMEOUT time.Duration = 5 * time.Second
	PORT    string        = ":8080"
)

func main() {
	// Copied from [there](https://github.com/gin-gonic/gin#graceful-restart-or-stop).

	// Create router.
	router := gin.Default()

	// Attach post handler.
	router.POST("/checkText", textChecker)

	srv := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		// Service connections.
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully
	// shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server was shutdowned.")
}

func textChecker(c *gin.Context) {
	// Fetching request data.
	dec := json.NewDecoder(c.Request.Body)
	req := Request{}
	err := dec.Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		c.Error(err).SetMeta("textChecker.Decode")
		return
	}

	bb := fasthttp.AcquireByteBuffer()
	searchTextBytes := []byte(req.SearchText)
	for _, site := range req.Site {
		// Get page's content.
		status, body, err := fasthttp.GetTimeout(bb.B, site, TIMEOUT)
		if err != nil || status != 200 {
			return
		}
		if bytes.Contains(body, searchTextBytes) {
			c.JSON(http.StatusOK, Response{site})
			fasthttp.ReleaseByteBuffer(bb)
			return
		}
	}
	fasthttp.ReleaseByteBuffer(bb)
	c.AbortWithStatus(http.StatusNoContent)
}

// Request contains sites and string
// that should be found on them.
type Request struct {
	Site       []string `json:"Site"`
	SearchText string   `json:"SearchText"`
}

// Response contains the first site
// where text from Request was found.
type Response struct {
	FoundAtSite string `json:"FoundAtSite"`
}
