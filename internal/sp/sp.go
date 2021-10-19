package apigw

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type sp struct {
	log *zap.Logger
}

type SP interface {
	Listen(port string) error
}

func (c *sp) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.log.Info("error in the REST handler", zap.Error(err))
}

func (c *sp) handle(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, res)
}

func (c *sp) Listen(port string) error {
	//
	//
	//

	// TODO: refactor and destructurize
	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, _ := url.Parse("http://localhost:8000")
	idpMetadataURL, _ := url.Parse("https://samltest.id/saml/idp")

	idpMetadata, err := samlsp.FetchMetadata(
		context.Background(),
		http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		IDPMetadata: idpMetadata,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		SignRequest: true,
	})
	if err != nil {
		panic(err) // TODO handle error
	}

	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(":8000", nil)

	//
	//
	//

	r := gin.Default()
	r.POST("/", c.handle)
	return r.Run(port)
}

func NewClient(log *zap.Logger) SP {
	return &sp{log}
}
