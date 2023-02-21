# Docs

## Swagger API spec

Use [swag](https://pkg.go.dev/github.com/go-openapi/swag) CLI to create API spec

```shell
cd starship/src/api-server/cmd
swag init -g ./api-server/http/http.go -d ../../
```

Start API Server, and visite the URL:
`http://<API-Server-address>/swagger/index.html`
