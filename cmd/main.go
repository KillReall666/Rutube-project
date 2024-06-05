package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/KillReall666/Rutube-project/internal/authentication"
	"github.com/KillReall666/Rutube-project/internal/config"
	"github.com/KillReall666/Rutube-project/internal/handlers/authorization"
	"github.com/KillReall666/Rutube-project/internal/handlers/getallusers"
	"github.com/KillReall666/Rutube-project/internal/handlers/registration"
	"github.com/KillReall666/Rutube-project/internal/handlers/subscribe"
	"github.com/KillReall666/Rutube-project/internal/handlers/unsubscribe"
	"github.com/KillReall666/Rutube-project/internal/interrogator"
	"github.com/KillReall666/Rutube-project/internal/logger"
	"github.com/KillReall666/Rutube-project/internal/service"
	"github.com/KillReall666/Rutube-project/internal/storage/postgres"
	"github.com/KillReall666/Rutube-project/internal/storage/redis"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic("failed to load cfg: " + err.Error())
	}

	log := logger.New()

	db, err := postgres.New(cfg.DefaultDBConnStr)
	if err != nil {
		log.LogFatal("db connection error:", err)
	}
	log.LogInfo("database connected")

	redisClient := redis.NewRedisClient(cfg.RedisAddress)
	pong, err := redisClient.Ping()
	if err != nil {
		log.LogFatal("redis connection error:", err)
	}
	log.LogInfo("connection to redis established:", pong)

	JWTMiddleware := authentication.JWTMiddleware{
		RedisClient: redisClient,
		Log:         log,
	}

	serv := service.New(cfg, log, db)

	interrog := interrogator.NewInterrogator(db, log)
	go func() {
		interrog.BirthDaysFinder()
	}()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware.JWTMiddleware())
		r.Post("/api/user/subscription", subscribe.NewSubscribeHandler(serv, log).Subscribe)
		r.Post("/api/user/unsubscription", unsubscribe.NewUnSubscribeHandler(serv, log).UnSubscribe)
		r.Post("/api/user/getdata", getallusers.NewGetAllUsersHandler(serv, log).GetAllUsers)
	})

	r.Post("/api/user/register", registration.NewRegistrationHandler(serv, redisClient, log).RegistrationHandler)
	r.Post("/api/user/login", authorization.NewAuthorizationHandler(serv, redisClient, log).AuthorizationHandler)

	log.LogInfo("starting server at localhost", cfg.Address)
	err = http.ListenAndServe(cfg.Address, r)
	if err != nil {
		log.LogFatal("server connection error:", err)
	}

}
