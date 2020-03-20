package main

import (
	"context"
	"flag"
	"fmt"
	"front-service/cmd/front/app"
	"front-service/cmd/front/app/rooms"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shuhrat-shokirov/core/pgk/core/auth"
	"github.com/shuhrat-shokirov/di/pkg/di"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/mux/pkg/mux"
	"log"
	"net"
	"net/http"
)

// flag - max priority, env - lower priority

var (
	host    = flag.String("host", "0.0.0.0", "Server host")
	port    = flag.String("port", "8888", "Server port")
	dsn    = flag.String("dsn", "postgres://user:pass@localhost:5403/rooms", "Server port")
	authUrl = flag.String("authUrl", "http://localhost:9999", "Auth Service URL")
)
type DSN string
func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)
	secret := jwt.Secret("secret")
	start(addr, secret, auth.Url(*authUrl),*dsn)
}

func start(addr string, secret jwt.Secret, authUrl auth.Url, dsn string) {
	container := di.NewContainer()
	err := container.Provide(
		app.NewServer,
		mux.NewExactMux,
		rooms.NewService,
		func() DSN { return DSN(dsn) },
		func(dsn DSN) *pgxpool.Pool {
			pool, err := pgxpool.Connect(context.Background(), string(dsn))
			if err != nil {
				panic(fmt.Errorf("can't create pool: %w", err))
			}
			return pool
		},
		func() jwt.Secret { return secret },
		func() auth.Url { return authUrl },
		auth.NewClient,
	)
	if err != nil {
		log.Fatal(err)
	}
	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	panic(http.ListenAndServe(addr, appServer))
}
