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
	getCSRFTokenRouter.HandleFunc("", handler.GetCSRFToken).Methods("GET")

	//
	signupRouter := r.PathPrefix("/signup").Subrouter()
	signupRouter.HandleFunc("", handler.CreateUser).Methods("POST")

	//
	signinRouter := r.PathPrefix("/signin").Subrouter()
	signinRouter.HandleFunc("", handler.Signin).Methods("POST")

	//
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.Use(middlewares.JWTValidate)
	usersRouter.HandleFunc("", handler.ListAllUsers).Methods("GET")

	//
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(middlewares.JWTValidate)
	{
		// get user by fullname
		userRouter.HandleFunc("/fullname/{fullname}", handler.SearchUserByFullname).Methods("GET")
		// get me
		userRouter.HandleFunc("/me", handler.GetUserDetailedInfo).Methods("GET")
		// delete user by account
		userRouter.HandleFunc("/account/{account}", handler.DeleteUser).Methods("DELETE")
		// update me
		userRouter.HandleFunc("/me", handler.UpdateUser).Methods("PATCH")
		// update specific user fullname
		userRouter.HandleFunc("/account/{account}", handler.UpdateSpecificUserFullname).Methods("PATCH")
	}

	// websocket
	wsRouter := r.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("/connection", handler.WebsocketConnection)
}
