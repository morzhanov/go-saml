package sp

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"

	"go.uber.org/zap"
)

type sp struct {
	log *zap.Logger
}

type SP interface {
	Listen() error
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", samlsp.AttributeFromContext(r.Context(), "cn"))
}

func (c *sp) Listen() error {
	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		return err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return err
	}

	rootURL, _ := url.Parse("http://localhost:8000")
	idpMetadataURL, _ := url.Parse("https://samltest.id/saml/idp")

	idpMetadata, err := samlsp.FetchMetadata(
		context.Background(),
		http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		return err
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		IDPMetadata: idpMetadata,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		SignRequest: true,
	})
	if err != nil {
		return err
	}

	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	return http.ListenAndServe(":8000", nil)
}

func NewClient(log *zap.Logger) SP {
	return &sp{log}
}
