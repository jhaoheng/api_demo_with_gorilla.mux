package route

import (
	"app/handler"
	"app/middlewares"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.Use(middlewares.ShowRequest)
	//
	getCSRFTokenRouter := r.PathPrefix("/csrf").Subrouter()
	getCSRFTokenRouter.HandleFunc("", handler.NewGetCSRFToken(nil)).Methods("GET")

	//
	signupRouter := r.PathPrefix("/signup").Subrouter()
	signupRouter.HandleFunc("", handler.NewCreateUser(nil)).Methods("POST")

	//
	signinRouter := r.PathPrefix("/signin").Subrouter()
	signinRouter.HandleFunc("", handler.NewSignin(nil)).Methods("POST")

	//
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.Use(middlewares.JWTValidate)
	usersRouter.HandleFunc("", handler.NewListAllUsers(nil)).Methods("GET")

	//
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(middlewares.JWTValidate)
	{
		// get user by fullname
		userRouter.HandleFunc("/fullname/{fullname}", handler.NewSearchUserByFullname(nil)).Methods("GET")
		// get me
		userRouter.HandleFunc("/me", handler.NewGetUserDetailed(nil)).Methods("GET")
		// delete user by account
		userRouter.HandleFunc("/account/{account}", handler.NewDeleteUser(nil)).Methods("DELETE")
		// update me
		userRouter.HandleFunc("/me", handler.NewUpdateUser(nil)).Methods("PATCH")
		// update specific user fullname
		userRouter.HandleFunc("/account/{account}", handler.NewUpdateUserFullname(nil)).Methods("PATCH")
	}

	// websocket
	wsRouter := r.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("/connection", handler.WebsocketConnection)
}
