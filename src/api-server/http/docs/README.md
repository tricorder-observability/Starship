# Docs

Swagger docs are served by HTTP serving handler, and is hosted together with
the HTTP server of which the HTTP APIs are described the Swagger docs..

Use [swag](https://pkg.go.dev/github.com/go-openapi/swag) CLI to generate
API spec:

```shell
# In the ToT of Starship repo
swag init -g src/api-server/http/http.go -o src/api-server/http/docs
```

To view API Server's Swagger APIs spec, start API Server, and visit the URL:
`http://<API-Server-address>/swagger/index.html`
