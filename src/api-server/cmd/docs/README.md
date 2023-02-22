# Docs

Swagger DOCs are HTTP serving handler that are hosted together with the HTTP
server.

[Wwag](https://pkg.go.dev/github.com/go-openapi/swag) CLI to create API spec:

```shell
cd src/api-server/cmd
swag init -g ./api-server/http/http.go -d ../../
```

To view API Server's Swagger APIs spec, start API Server, and visite the URL:
`http://<API-Server-address>/swagger/index.html`
