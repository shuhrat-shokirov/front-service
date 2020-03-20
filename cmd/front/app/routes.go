package app

import (
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/authenticated"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/jwt"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/logger"
	"github.com/shuhrat-shokirov/mux/pkg/mux/middleware/unauthenticated"
	"reflect"
)

var (
	Root   = "/"
	Login  = "/login"
	Logout = "/logout"
	Posts  = "/posts"
	Post     = "/posts/%s"
	PostEdit = "/posts/%s/edit"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwt.SourceCookie, true, Logout, reflect.TypeOf((*Payload)(nil)).Elem(), s.secret)
	authMW := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)
	unAuthMW := unauthenticated.Unauthenticated(jwt.IsContextNonEmpty, true, "/me")
	s.router.GET("/", s.handleFrontPage(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/me", s.handleFrontPageForAuth(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET(Login, s.handleLoginPage(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Login, s.handleLogin(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET(Logout, s.handleLogout(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Logout, s.handleLogout(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST("/register", s.handleRegister(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/register", s.handleRegisterGet(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET(PostEdit, s.handlePostEditPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(PostEdit, s.handlePostEdit(), authMW, jwtMW, logger.Logger("HTTP"))
}
