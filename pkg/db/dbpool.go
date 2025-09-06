package dbpool

import (
	"context"
	"fmt"
	"profile-service/pkg/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func BuildDSN(cfg *config.Config) string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
	return dsn
}

func InitDBPool(dsn string, cfg *config.Config, log *logrus.Logger, ctx context.Context) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn string is empty")
	}

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.WithError(err).Error("failed parse dsn db config")
		return nil, err
	}

	conf.MaxConnIdleTime = 30 * time.Second
	conf.MaxConnLifetime = 10 * time.Minute
	conf.MaxConns = 10
	conf.MinConns = 2
	conf.HealthCheckPeriod = 5 * time.Minute

	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var pool *pgxpool.Pool
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(connCtx, conf)
		if err != nil {
			log.WithError(err).Errorf("failed create new pool, attempts: %d", i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			log.WithError(err).Warnf("failed ping db")
			pool.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		log.Infof("Pool connection succeeded after %d attempts", i+1)
		return pool, nil
	}

	log.WithError(err).Errorf("failed create db pool after 5 tries: %v", err)
	return nil, err
}
