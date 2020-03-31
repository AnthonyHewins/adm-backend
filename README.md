ADM backend
===========

1. [Development](#development)
2. [Build](#build)
3. [Production](#production)

### Development

Uses `gin-gonic`.

Tip: run these all at once with `Ctrl+X, Ctrl+E` in your terminal

1. `git clone https://github.com/AnthonyHewins/adm-backend`
2. `go mod download`
3. `go run main.go`
4. `curl localhost:8080` and you should see a `404`

### Build

Uses multi-stage builds in docker. End image is a scratch image with only the binary

1. `docker build -t ahewins/adm-backend:$TAG .`
2. `docker push ahewins/adm-backend:$TAG`

Compressed into one command:

``` sh
TAG=something;docker build -t ahewins/adm-backend:$TAG && docker push ahewins/adm-backend:$TAG
```

### Production

1. `docker run -p 8080:8080 -d ahewins/adm-backend`
