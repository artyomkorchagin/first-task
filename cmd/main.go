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

	"github.com/IBM/sarama"
	"github.com/artyomkorchagin/first-task/internal/config"
	"github.com/artyomkorchagin/first-task/internal/infrastructure"
	"github.com/artyomkorchagin/first-task/internal/logger"
	orderpostgresql "github.com/artyomkorchagin/first-task/internal/repository/postgres/order"
	"github.com/artyomkorchagin/first-task/internal/router"
	orderservice "github.com/artyomkorchagin/first-task/internal/service"
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
	var err error

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	var zapLogger *zap.Logger

	if cfg.LogMode == "DEV" {
		zapLogger, err = logger.NewDevelopmentLogger()
	} else {
		zapLogger, err = logger.NewLogger()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	zapLogger.Info("Starting application")
	zapLogger.Info("Connecting to database")
	db, err := sql.Open("pgx", cfg.GetDSN())
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	zapLogger.Info("Pinging database")
	if err := db.Ping(); err != nil {
		zapLogger.Fatal("Failed to ping database", zap.Error(err))
	}
	zapLogger.Info("Connected to database")

	zapLogger.Info("Running migrations")
	if err := orderpostgresql.RunMigrations(db); err != nil {
		zapLogger.Fatal("Failed to run up migration", zap.Error(err))
	}
	zapLogger.Info("Succesfully ran up migration")

	zapLogger.Info("Connecting to redis")
	rAddr := cfg.Redis.Host + cfg.Redis.Port
	rdb, err := infrastructure.NewRedisClient(rAddr, cfg.Redis.Password)
	if err != nil {
		zapLogger.Fatal("Failed to connect to redis", zap.Error(err))
	}

	zapLogger.Info("Connected to redis")

	worker, err := infrastructure.ConnectConsumer(cfg.Kafka)
	if err != nil {
		zapLogger.Fatal("Failed to connect consumer to kafka", zap.Error(err))
	}

	consumer, err := worker.ConsumePartition(cfg.Kafka.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		zapLogger.Fatal("Failed to consume partition", zap.Error(err))
	}

	walletRepo := orderpostgresql.NewRepository(db)
	walletSvc := orderservice.NewService(walletRepo, rdb, zapLogger)

	handler := router.NewHandler(walletSvc, zapLogger)
	r := handler.InitRouter()

	port := cfg.Server.Port
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		zapLogger.Info("Starting Kafka consumer...")
		defer zapLogger.Info("Kafka consumer stopped")
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-consumer.Errors():
				zapLogger.Error("Recieved error from kafka", zap.Error(err))
			case msg := <-consumer.Messages():
				zapLogger.Info("Got message", zap.String("topic", msg.Topic),
					zap.Int32("partition", msg.Partition),
					zap.Int64("offset", msg.Offset))
				handler.CreateOrderKafka(ctx, msg.Value)
			}
		}
	}()

	go func() {
		zapLogger.Info("Server starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	<-quit
	zapLogger.Info("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zapLogger.Error("Server shutdown failed", zap.Error(err))
	}
	zapLogger.Info("Server exited")

	if err := consumer.Close(); err != nil {
		zapLogger.Error("Failed to close Kafka consumer", zap.Error(err))
	}

	if err := worker.Close(); err != nil {
		zapLogger.Error("Failed to close Kafka worker", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		zapLogger.Error("Error closing database connection", zap.Error(err))
	}

	zapLogger.Info("Shutdown completed")
}
