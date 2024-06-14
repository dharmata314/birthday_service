package main

import (
	"birthday-service/internal/config"
	"birthday-service/internal/database"
	database3 "birthday-service/internal/database/emp_repo"
	database2 "birthday-service/internal/database/subs_repo"
	database4 "birthday-service/internal/database/user_repo"
	errMsg "birthday-service/internal/err"
	handlers2 "birthday-service/internal/handlers/emp"
	handlers3 "birthday-service/internal/handlers/subs"
	handlers "birthday-service/internal/handlers/user"
	notification "birthday-service/internal/notification"
	"birthday-service/jwt"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const notificationFrequency = 1

func main() {
	cfg := config.MustLoad()
	log := setupLogger()
	fmt.Println(cfg)
	log.Debug("debug messages are active")
	pg, err := connectToPostgres(cfg, log)
	if err != nil {
		log.Error("failed to create postgres db", errMsg.Err(err))
		os.Exit(1)
	}
	fmt.Println("connecting to postgres...")
	defer pg.Close()
	if pg == nil {
		log.Error("failed to connect to postgres")
		os.Exit(1)
	}
	if err := pg.Ping(context.Background()); err != nil {
		log.Error("failed to ping postgres db", errMsg.Err(err))
		os.Exit(1)
	} else {
		log.Info("postgres db connected successfully")
	}

	log.Info("application started")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	cfgSMTP := &config.ConfigSMTP{
		SMTPHost:     "smtp.yandex.ru",
		SMTPPort:     587,
		SMTPUsername: "",
		SMTPPassword: "",
	}

	empRepository := database3.NewEmployeeRepository(pg.Db, log)
	subsRepository := database2.NewSubsRepository(pg.Db, log)
	userRepository := database4.NewUserRepository(pg.Db, log)
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, log)

	router.Post("/users/new", handlers.New(log, userRepository))
	router.Post("/login", handlers.LoginFunc(log, userRepository, jwtManager))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Delete("/users/{id}", handlers.DeleteUserHandler(log, userRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Patch("/users/{id}", handlers.NewUpdateUserHandler(userRepository, log))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Post("/emp", handlers2.New(log, empRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Delete("/emp/{id}", handlers2.DeleteEmpHandler(log, empRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Get("/employees", handlers2.ListAllEmployees(log, empRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Post("/subs", handlers3.New(log, subsRepository))

	router.With(func(next http.Handler) http.Handler {
		return jwt.TokenAuthMiddleware(jwtManager, next)
	}).Delete("/subs/{id}", handlers3.DeleteSub(log, subsRepository))

	log.Info("starting server", slog.String("addr", cfg.HTTPServer.Addr))
	server := &http.Server{
		Addr:              cfg.HTTPServer.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server", errMsg.Err(err))
		}
	}()

	for {
		ctx := context.Background()
		notification.SendBirthdayNotifications(ctx, subsRepository, empRepository, cfgSMTP, log)
		time.Sleep(notificationFrequency * time.Minute)
	}

}

func setupLogger() *slog.Logger {
	var log *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

func connectToPostgres(cfg *config.Config, log *slog.Logger) (*database.Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	pg, err := database.NewPG(context.Background(), connString, log, cfg)
	return pg, err
}
