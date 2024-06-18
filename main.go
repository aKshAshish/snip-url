package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var l = log.New(os.Stdout, "snip-url:", log.LstdFlags)

type ShortUrl struct {
	SUrl string `uri:"sUrl" binding:"required"`
}

type Url struct {
	Url  string `json:"url" binding:"required"`
	SUrl string `json:"sUrl"`
}

func getUrl(c *gin.Context) {
	l.Println("Getting the url")
	var sUrl ShortUrl
	if err := c.ShouldBindUri(&sUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "No valid short-url present with the request"})
		return
	}
	l.Printf("Short Url: %s\n", sUrl.SUrl)
}

func addUrl(c *gin.Context) {
	l.Println("Adding the url")
	var url Url
	if err := c.ShouldBindBodyWithJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid body."})
		return
	}
	l.Printf("Request: %+v\n", url)
}

func main() {

	router := gin.Default()

	router.GET("/:sUrl", getUrl)
	router.POST("/url", addUrl)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Println("Starting server on :9090")
		err := s.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			l.Printf("Error while starting the server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm and iterrupt signals to shutdown gracefully
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	sig := <-ch

	l.Printf("Got Signal: %s\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		l.Fatal("Error shutting the server down:", err)
	}

	l.Println("Sever exiting.")
}
