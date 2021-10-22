package idp

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"

	"github.com/crewjam/saml/logger"
	"github.com/crewjam/saml/samlidp"
	"github.com/zenazn/goji"
	"golang.org/x/crypto/bcrypt"
)

type idp struct {
	log     logger.Interface
	store   samlidp.Store
	baseURL string
	cert    *x509.Certificate
	key     *rsa.PrivateKey
}

type IDP interface {
	Listen() error
}

func (i *idp) Listen() error {
	baseURL, err := url.Parse(i.baseURL)
	if err != nil {
		return fmt.Errorf("cannot parse base URL: %v", err)
	}

	idpServer, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Key:         i.key,
		Logger:      i.log,
		Certificate: i.cert,
		Store:       i.store,
	})
	if err != nil {
		i.log.Fatalf("%s", err)
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
		i.log.Fatalf("%s", err)
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
		i.log.Fatalf("%s", err)
	}

	goji.Handle("/*", idpServer)
	goji.Serve()
	return nil
}

func NewServer(log logger.Interface, store samlidp.Store, baseURL string, key string, cert string) (IDP, error) {
	b, _ := pem.Decode([]byte(key))
	k, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	b, _ = pem.Decode([]byte(cert))
	if err != nil {
		return nil, err
	}
	c, _ := x509.ParseCertificate(b.Bytes)

	return &idp{log, store, baseURL, c, k}, nil
}
