package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artyomkorchagin/first-task/internal/config"
	"github.com/artyomkorchagin/first-task/internal/logger"
	orderpostgresql "github.com/artyomkorchagin/first-task/internal/repository/postgres/order"
	"github.com/artyomkorchagin/first-task/internal/router"
	orderservice "github.com/artyomkorchagin/first-task/internal/service"
	"github.com/artyomkorchagin/first-task/pkg/helpers"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//	@title			WB First Task
//	@version		1.0

//	@contact.name	Artyom Korchagin
//	@contact.email	artyomkorchagin333@gmail.com

//	@host		localhost:3000
//	@BasePath	/

func main() {

	var zapLogger *zap.Logger
	var err error

	if helpers.GetEnv("ENV", "DEV") == "DEV" {
		zapLogger, err = logger.NewDevelopmentLogger()
	} else {
		zapLogger, err = logger.NewLogger()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting application")

	db, err := sql.Open("pgx", config.GetDSN())
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		zapLogger.Fatal("Failed to ping database", zap.Error(err))
	}
	zapLogger.Info("Connected to database")

	if err := orderpostgresql.RunMigrations(db); err != nil {
		zapLogger.Fatal("Failed to run up migration", zap.Error(err))
	}
	zapLogger.Info("Succesfully ran up migration")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: helpers.GetEnv("REDIS_PASSOWORD", ""),
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		zapLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	repo := orderpostgresql.NewRepository(db)
	service := orderservice.NewService(repo)

	handler := router.NewHandler(service, rdb, zapLogger)
	r := handler.InitRouter()

	srv := &http.Server{
		Addr:    helpers.GetEnv("SERVER_HOST", "") + ":" + helpers.GetEnv("SERVER_PORT", "3000"),
		Handler: r,
	}

	go func() {
		zapLogger.Info("Server starting", zap.String("port", helpers.GetEnv("SERVER_PORT", "3000")))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Server shutdown failed", zap.Error(err))
	}

	zapLogger.Info("Server exited")

	if err := db.Close(); err != nil {
		zapLogger.Error("Error closing database connection", zap.Error(err))
	}
	zapLogger.Info("Closed database connection")

	zapLogger.Info("Program exited")
}
