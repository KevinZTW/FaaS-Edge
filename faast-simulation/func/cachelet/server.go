package cachelet

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PeerConnectRequest struct {
	endpoint string `json:"endpoint"`
}

type PeerConnectResponse struct {
	message string `json:"message"`
}

func (c *Cachelet) initServer() error {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cachelets", func(ctx *gin.Context) {
		key := ctx.Query("key")
		value := c.Get(key)

		ctx.JSON(http.StatusOK, gin.H{
			"key":   key,
			"value": value,
		})
	})

	r.POST("/peers", func(ctx *gin.Context) {
		req := PeerConnectRequest{}
		ctx.Bind(&req)
		if err := c.handleConnection(req.endpoint); err != nil {
			ctx.Error(err)
		} else {
			ctx.JSON(http.StatusOK, gin.H{})
		}
	})

	go func() {
		r.Run(":3038")
	}()
	return nil
}
