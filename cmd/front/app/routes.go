package app

import (
	"front-service/pkg/mux/middleware/authenticated"
	"front-service/pkg/mux/middleware/jwt"
	"front-service/pkg/mux/middleware/logger"
	"front-service/pkg/mux/middleware/unauthenticated"
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
	s.router.POST("/", s.handleFrontPage(), unAuthMW, jwtMW, logger.Logger("HTTP"))

	s.router.POST("/me", s.handleFrontPageForAuth(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/me", s.handleFrontPageForAuth(), authMW, jwtMW, logger.Logger("HTTP"))

	s.router.GET(Login, s.handleLoginPage(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Login, s.handleLogin(), unAuthMW, jwtMW, logger.Logger("HTTP"))

	s.router.GET(Logout, s.handleLogout(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Logout, s.handleLogout(), authMW, jwtMW, logger.Logger("HTTP"))

	s.router.POST("/register", s.handleRegister(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/register", s.handleRegisterGet(), unAuthMW, jwtMW, logger.Logger("HTTP"))

	s.router.GET(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))
	s.router.POST(Posts, s.handlePostsPage(), authMW, jwtMW, logger.Logger("HTTP"))

	s.router.POST("/rooms/0", s.handleAddNewRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/rooms/0", s.handleAddNewRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))

	s.router.GET("/rooms/history/{id}", s.handleHistoryRoom(), jwtMW, logger.Logger("HTTP"))
	s.router.POST("/rooms/history/{id}", s.handleHistoryRoom(), jwtMW, logger.Logger("HTTP"))

	s.router.GET("/rooms", s.handleRoom(), jwtMW, logger.Logger("HTTP"))
	s.router.POST("/rooms", s.handleRoom(), jwtMW, logger.Logger("HTTP"))

	s.router.GET("/history/rooms", s.handleOpenRooms(), jwtMW, logger.Logger("HTTP"))
	s.router.POST("/history/rooms", s.handleOpenRooms(), jwtMW, logger.Logger("HTTP"))

	s.router.POST("/rooms/history/0", s.handleAddNewHistoryRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/rooms/history/0", s.handleAddNewHistoryRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))

	s.router.POST("/rooms/history/result/{id}", s.handleAddNewResultHistoryRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	s.router.GET("/rooms/history/result/{id}", s.handleAddNewResultHistoryRoom(), unAuthMW, jwtMW, logger.Logger("HTTP"))

}
