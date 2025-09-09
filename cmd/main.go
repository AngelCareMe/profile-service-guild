package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"profile-service/internal/adapter/blizzard"
	"profile-service/internal/adapter/database"
	"profile-service/internal/handler"
	"profile-service/internal/usecase"
	"profile-service/pkg/config"
	dbpool "profile-service/pkg/db"
	logger "profile-service/pkg/log"
	"syscall"
	"time"

	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("load config error")
	}

	log := logger.InitLogger(cfg)

	ctx := context.Background()
	dsn := dbpool.BuildDSN(cfg)
	pool, err := dbpool.InitDBPool(dsn, cfg, log, ctx)
	if err != nil {
		log.Fatalf("Failed init db pool")
	}
	defer pool.Close()

	if err := runMigrations(pool, log); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	blizzAd := blizzard.NewBlizzardRepository(log)
	dbAd := database.NewPostgresRepository(pool, log)

	profileUc := usecase.NewProfileUsecase(dbAd, blizzAd, log)

	profileHandl := handler.NewProfileHandler(blizzAd, profileUc, log)

	router := gin.Default()
	handler.SetupRoutes(router, profileHandl, cfg, log)

	log.Infof("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed start server: %v", err)
		}
	}()

	log.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Server shut down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed shutdown server: %v", err)
	}

	log.Info("Server exited")
}

func runMigrations(pool *pgxpool.Pool, log *logrus.Logger) error {
	start := time.Now()
	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Warn("closing migrate db failed")
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.WithError(srcErr).Warn("migrate source close failed")
		}
		if dbErr != nil {
			log.WithError(dbErr).Warn("migrate db close failed")
		}
	}()

	log.Info("Applying migrations...")

	var migErr error
	err = retry.Do(
		func() error {
			_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			migErr = m.Up()
			if migErr != nil && !errors.Is(migErr, migrate.ErrNoChange) {
				return fmt.Errorf("migrate up failed: %w", migErr)
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)
	if err != nil {
		return err
	}

	if errors.Is(migErr, migrate.ErrNoChange) {
		log.Info("No migrations to apply")
	}

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("get migrate version: %w", err)
	}
	log.WithFields(logrus.Fields{
		"duration": time.Since(start),
		"version":  version,
		"dirty":    dirty,
	}).Info("Migrations succeeded")

	return nil
}
