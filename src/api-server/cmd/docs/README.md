# API Server http document

## create swagger
- use swag cli create docs[swag](https://pkg.go.dev/github.com/go-openapi/swag)
```shell
cd starship/src/api-server/cmd
swag init -g ./api-server/http/http.go -d ../../
```

## Visit document
http://host:port/swagger/index.html