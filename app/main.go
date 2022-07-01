package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"api_demo_with_gorilla.mux/app/config"
	"api_demo_with_gorilla.mux/app/middlewares"
	"api_demo_with_gorilla.mux/app/models"
	"api_demo_with_gorilla.mux/app/modules"
	"api_demo_with_gorilla.mux/app/route"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func init() {
	env := os.Getenv("env")
	if env == "prod" {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	c := config.NewConfig(env)
	models.NewDBMySQL(models.DBSet{
		Host:    c.DB_HOST,
		User:    c.DB_USERNAME,
		Pass:    c.DB_PASSWORD,
		DBName:  c.DB_NAME,
		IsDebug: true,
	})
	//
	modules.InitValidate()
}

func main() {
	Run()
}

func Run() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	//
	r := mux.NewRouter()
	r.Use(middlewares.ShowRequest)
	r.Use(set_CSRF()) // CSRF protection
	route.RegisterRoutes(r)
	route.WalkingRoute(r)

	//
	srv := &http.Server{
		Addr: ":8080",
		// to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      time.Second * 15,
		ReadTimeout:       time.Second * 15,
		IdleTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 15,
		Handler:           r,
	}

	logrus.Info("api start")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	//
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	//
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	logrus.Info("shutting down")
	os.Exit(0)
}

func set_CSRF() func(http.Handler) http.Handler {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return csrf.Protect(
		key,
		csrf.Secure(config.CFG.CSRFTOKEN_ONLY_HTTPS),
		csrf.TrustedOrigins([]string{"*"}), // ex: localhost, www.google.com
		csrf.Domain("localhost"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := fmt.Errorf("forbidden - CSRF token invalid")
			status := http.StatusForbidden
			modules.NewResp(w, r).Set(modules.RespContect{Stutus: status, Error: err})
		})),
	)
}
