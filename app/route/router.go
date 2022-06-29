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
	getCSRFTokenRouter.HandleFunc("", handler.GetCSRFTokenHandler).Methods("GET")

	//
	signupRouter := r.PathPrefix("/signup").Subrouter()
	signupRouter.HandleFunc("", handler.CreateUserHandler).Methods("POST")

	//
	signinRouter := r.PathPrefix("/signin").Subrouter()
	signinRouter.HandleFunc("", handler.NewSignin(nil)).Methods("POST")

	//
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.Use(middlewares.JWTValidate)
	usersRouter.HandleFunc("", handler.ListAllUsersHandler).Methods("GET")

	//
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(middlewares.JWTValidate)
	{
		// get user by fullname
		userRouter.HandleFunc("/fullname/{fullname}", handler.NewSearchUserByFullname(nil)).Methods("GET")
		// get me
		userRouter.HandleFunc("/me", handler.GetUserDetailedHandler).Methods("GET")
		// delete user by account
		userRouter.HandleFunc("/account/{account}", handler.DeleteUserHandler).Methods("DELETE")
		// update me
		userRouter.HandleFunc("/me", handler.UpdateUserHandler).Methods("PATCH")
		// update specific user fullname
		userRouter.HandleFunc("/account/{account}", handler.NewUpdateUserFullname(nil)).Methods("PATCH")
	}

	// websocket
	wsRouter := r.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("/connection", handler.WebsocketConnection)
}
