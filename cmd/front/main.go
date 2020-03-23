package main

import (
	"flag"
	"front-service/cmd/front/app"
	"github.com/shuhrat-shokirov/core/pgk/core/auth"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/new-mux/pkg/mux"
	"net"
	"net/http"
)

// flag - max priority, env - lower priority

var (
	host    = flag.String("host", "0.0.0.0", "Server host")
	port    = flag.String("port", "8888", "Server port")
	authUrl = flag.String("authUrl", "http://localhost:9999", "Auth Service URL")
)
type DSN string
func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	secret := jwt.Secret("secret")
	start(addr, secret, auth.Url(*authUrl))
}

func start(addr string, secret jwt.Secret, authUrl auth.Url) {
	exactMux := mux.NewExactMux()
	client := auth.NewClient(authUrl)
	server := app.NewServer(exactMux, secret, client)
	server.Start()
	panic(http.ListenAndServe(addr, server))
}
