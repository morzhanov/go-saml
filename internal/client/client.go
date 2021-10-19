package apigw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type client struct {
	log *zap.Logger
}

type Client interface {
	Listen(port string) error
}

func (c *client) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.log.Info("error in the REST handler", zap.Error(err))
}

func (c *client) handle(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, res)
}

func (c *client) Listen(port string) error {
	r := gin.Default()
	r.POST("/", c.handle)
	return r.Run(port)
}

func NewClient(log *zap.Logger) Client {
	return &client{log}
}
