package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/shuhrat-shokirov/core/pgk/core/auth"
	"github.com/shuhrat-shokirov/core/pgk/core/utils"
	"github.com/shuhrat-shokirov/jwt/pkg/cmd"
	"github.com/shuhrat-shokirov/mux/pkg/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Server struct {
	router     *mux.ExactMux
	secret     jwt.Secret
	authClient *auth.Client
}

func NewServer(router *mux.ExactMux, secret jwt.Secret, authClient *auth.Client) *Server {
	return &Server{router: router, secret: secret, authClient: authClient}
}
func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleFrontPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)

	tpl, err = template.ParseFiles(filepath.Join("web/templates", "index.gohtml"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err = tpl.Execute(writer, struct {}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
		//log.Print(rooms)
	}
}

func (s *Server) handleFrontPageForAuth() http.HandlerFunc {
	// executes in one goroutine
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates/users", "user.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err = tpl.Execute(writer, struct {}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleLoginPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "login.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "login.gohtml"))
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Printf("error while parse login form: %v", err)
			return
		}

		login := request.PostFormValue("login")
		if login == "" {
			log.Print("login can't be empty")
			return
		}
		password := request.PostFormValue("password")
		if password == "" {
			//
			log.Print("password can't be empty")
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		request = request.WithContext(ctx)

		token, err := s.authClient.Login(request.Context(), login, password)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				log.Print("auth service didn't response in given time")
				log.Print("another err") // parse it
			case errors.Is(err, context.Canceled):
				log.Print("auth service didn't response in given time")
				log.Print("another err") // parse it
			case errors.Is(err, auth.ErrResponse):
				var typedErr *auth.ErrorResponse
				ok := errors.As(err, &typedErr)
				if ok {
					tplData := struct {
						Err string
					}{
						Err: "",
					}
					if utils.StringInSlice("err.password_mismatch", typedErr.Errors) {
						tplData.Err = "err.password_mismatch"
					}
					err := tpl.Execute(writer, tplData)
					if err != nil {
						log.Print(err)
					}
				}
			}
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
		http.Redirect(writer, request, Posts, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handlePostsPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handlePostEditPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handlePostEdit() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cookie := &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(writer, cookie)
	}
}

func (s *Server) handleRegisterGet() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(filepath.Join("web/templates", "register.gohtml"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, nil)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Print(err)
		}
		name := request.PostForm.Get("name")
		log.Print(name)
		login := request.PostFormValue("login")
		password := request.PostFormValue("password")
		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		err = s.authClient.Register(ctx, name, login, password)
		log.Print(err)
		if err != nil {
			if err == auth.ErrAddNewUser {
				writer.Write([]byte("Пользователь с таким логином Существует"))
				return
			} else {
				log.Printf("что-то не то, %v", err)
			}
		} else {
			writer.Write([]byte("Пользователь успешно зарегистрирован!"))
			return
		}

	}
}


func (s *Server) handleAddNewRoom() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			log.Print(err)
		}
		book := auth.Rooms{}
		all, err := ioutil.ReadAll(request.Body)
		err = json.Unmarshal(all, &book)
		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		err = s.authClient.NewRoom(ctx, book)
		log.Print(err)
		if err != nil {
				log.Printf("что-то не то, %v", err)
			} else {
			writer.Write([]byte("Комната успешно дабавлена!"))
			return
		}

	}
}

//func (s *Server) handleHistoryRoom() http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		idFromCTX, ok := mux.FromContext(request.Context(), "id")
//		if !ok {
//			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		id, err := strconv.Atoi(idFromCTX)
//		if err != nil {
//			http.Error(writer,http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		err = request.ParseForm()
//		if err != nil {
//			log.Print(err)
//		}
//		book := auth.RoomsHistory{}
//		all, err := ioutil.ReadAll(request.Body)
//		err = json.Unmarshal(all, &book)
//		ctx, _ := context.WithTimeout(request.Context(), time.Second)
//		err = s.authClient.HistoryRoom(ctx, int64(id))
//		log.Print(err)
//		if err != nil {
//				log.Printf("что-то не то, %v", err)
//			} else {
//			writer.Write([]byte("История комнат под id!"))
//			return
//		}
//
//	}
//}
