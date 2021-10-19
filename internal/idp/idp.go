package apigw

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/logger"
	"github.com/crewjam/saml/samlidp"
	"github.com/gin-gonic/gin"
	"github.com/zenazn/goji"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type idp struct {
	log *zap.Logger
}

type IDP interface {
	Listen(port string) error
}

func (c *idp) handleHttpErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
	c.log.Info("error in the REST handler", zap.Error(err))
}

func (c *idp) handle(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, res)
}

func (c *idp) Listen(port string) error {
	//
	//
	//

	// TODO: refactor and destructurize
	logr := logger.DefaultLogger
	baseURLstr := flag.String("idp", "", "The URL to the IDP")
	flag.Parse()

	baseURL, err := url.Parse(*baseURLstr)
	if err != nil {
		logr.Fatalf("cannot parse base URL: %v", err)
	}

	idpServer, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Key:         key,
		Logger:      logr,
		Certificate: cert,
		Store:       &samlidp.MemoryStore{},
	})
	if err != nil {
		logr.Fatalf("%s", err)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hunter2"), bcrypt.DefaultCost)
	err = idpServer.Store.Put("/users/alice", samlidp.User{Name: "alice",
		HashedPassword: hashedPassword,
		Groups:         []string{"Administrators", "Users"},
		Email:          "alice@example.com",
		CommonName:     "Alice Smith",
		Surname:        "Smith",
		GivenName:      "Alice",
	})
	if err != nil {
		logr.Fatalf("%s", err)
	}

	err = idpServer.Store.Put("/users/bob", samlidp.User{
		Name:           "bob",
		HashedPassword: hashedPassword,
		Groups:         []string{"Users"},
		Email:          "bob@example.com",
		CommonName:     "Bob Smith",
		Surname:        "Smith",
		GivenName:      "Bob",
	})
	if err != nil {
		logr.Fatalf("%s", err)
	}

	goji.Handle("/*", idpServer)
	goji.Serve()

	//
	//
	//
	//

	r := gin.Default()
	r.POST("/", c.handle)
	return r.Run(port)
}

func NewServer(log *zap.Logger) IDP {
	return &idp{log}
}
