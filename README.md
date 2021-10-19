# go-saml

Go SAML example based on <a href="https://github.com/crewjam/saml/tree/main/example">go-saml example</a>

- `IDP` - Identity Provider service
- `SP` - Service Provider service

## Running

- Run each service separately and perform `POST /auth` request to SP
- SP should issue SAML token
- Now client able to perform `GET /info` endpoint request with SAML token

```bash
go run ./cmd/idp.go
go run ./cmd/sp.go

curl -d "TODO: add data" -X POST http://localhost:8080/auth

curl http://localhost:8080/info TODO: add header
```
