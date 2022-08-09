package database

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"loyalty-service/internal/logger"
)

func NewDatabase(DSN string, logger *logger.Logger) *sqlx.DB {
	db, err := sqlx.Connect("pgx", DSN)
	if err != nil {
		log.Fatalln(err)
	}

	migrations := NewMigrationsHandler(db, logger)
	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("migrations error", zap.Error(err))
	}
	logger.Info("migrations up successfully")

	return db
}

func NewMigrationsHandler(db *sqlx.DB, logger *logger.Logger) *migrate.Migrate {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatal("migrations error", zap.Error(err))
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"postgres", driver)
	if err != nil {
		logger.Fatal("migrations error", zap.Error(err))
	}

	return migrations
}
