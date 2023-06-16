package main

import (
	"func/app"
	"func/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	var peerAddress string
	var port string

	// used in non-container dev
	if len(os.Args) == 3 {
		peerAddress = os.Args[1]
		port = os.Args[2]
	}

	if port == "" {
		port = "3038"
	}
	st := store.New(peerAddress)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/store", func(c *gin.Context) {
		key := c.Query("key")
		value := st.PeersGet(key)

		c.JSON(http.StatusOK, gin.H{
			"key":   key,
			"value": value,
		})

	})

	go func() {
		workload := app.NewBasicWorkload(peerAddress)
		workload.Run()
		workload.ReportOutCome()
	}()

	r.Run(":" + port)
}
