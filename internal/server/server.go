package apigw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	log *zap.Logger
}

type Server interface {
	Listen(port string) error
}

func (c *server) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.log.Info("error in the REST handler", zap.Error(err))
}

func (c *server) handle(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, res)
}

func (c *server) Listen(port string) error {
	r := gin.Default()
	r.POST("/", c.handle)
	return r.Run(port)
}

func NewServer(log *zap.Logger) Server {
	return &server{log}
}
