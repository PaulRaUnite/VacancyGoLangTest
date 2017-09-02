// +build go1.8
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/PaulRaUnite/VacancyGoLangTest/util"
	"github.com/gin-gonic/gin"
)

func main() {
	// Copied from [there](https://github.com/gin-gonic/gin#graceful-restart-or-stop).

	// Create router.
	router := gin.Default()

	// Attach post handler.
	router.POST("/checkText", textChecker)

	srv := &http.Server{
		Addr:    ":8080",
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
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exist")
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

	// Create channel for workers
	// of amount of pages.
	// Buffer is important, because
	// all of workers can return
	// urls, so channel and WaitGroup
	// can make deadlock without it.
	outcome := make(chan string, len(req.Site))
	wg := sync.WaitGroup{}
	for _, site := range req.Site {
		// Start concurrent workers.
		wg.Add(1)
		go scan(site, req.SearchText, outcome, &wg)
	}

	// Wait for workers.
	wg.Wait()
	close(outcome)

	// Get the first site where the text was found.
	page, ok := <-outcome
	if ok {
		c.JSON(http.StatusOK, Response{page})
	} else {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

var timeoutClient = &http.Client{
	Timeout: 10 * time.Second,
}

func scan(url, text string, out chan<- string, group *sync.WaitGroup) {
	// Done work in all cases.
	defer func() {
		group.Done()
	}()

	// Get page's content.
	resp, err := timeoutClient.Get(url)
	if err != nil {
		return
	}

	// Searching text on the url "streamly".
	// Not sure(it means I have no benchmarks)
	// that it is faster than naїve
	// searching but it doesn't require
	// additional memory for byte slice,
	// so maybe it is.
	if util.Search(resp.Body, text) {
		// Return url.
		out <- url
	}
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
