package main

import (
	"context"
	"crypto/rand"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"api_demo_with_gorilla.mux/app/config"
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
	r.Use(setCSRF()) // CSRF protection
	route.RegisterRoutes(r)
	route.WalkingRoute(r)

	//
	srv := &http.Server{
		Addr: ":8080",
		// to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
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
	srv.Shutdown(ctx)
	logrus.Info("shutting down")
	os.Exit(0)
}

func setCSRF() func(http.Handler) http.Handler {
	key := make([]byte, 32)
	rand.Read(key)
	return csrf.Protect(
		key,
		csrf.Secure(config.CFG.CSRFTOKEN_ONLY_HTTPS),
		// csrf.TrustedOrigins([]string{""}),
		csrf.Domain("localhost"),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)
}
