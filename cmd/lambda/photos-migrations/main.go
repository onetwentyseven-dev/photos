package main

import (
	"context"
	"database/sql"
	"photos"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrateDriver "github.com/golang-migrate/migrate/v4/database/mysql"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

type app struct {
	logger *logrus.Logger
	db     *sql.DB
}

type payload struct {
	Operation string `json:"operation"`
}

func (a *app) handle(ctx context.Context, operation *payload) error {

	a.logger.WithField("operation", operation.Operation).Info("running migrations")

	migrateDriver, err := mysqlMigrateDriver.WithInstance(a.db, &mysqlMigrateDriver.Config{
		MigrationsTable:  "migrations",
		DatabaseName:     appConfig.DB.Name,
		StatementTimeout: time.Second * 5,
	})
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize migration driver")
	}

	migrateSource, err := photos.MigrationFS()
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize migration source")
	}

	m, err := migrate.NewWithInstance("iofs", migrateSource, "mysql", migrateDriver)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize migration instance")
	}

	a.logger.WithField("operation", operation.Operation).Info("migration loaded, starting execution")

	switch operation.Operation {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			logger.WithError(err).Fatal("failed to run migrations")
		}
	case "down":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			logger.WithError(err).Fatal("failed to run migrations")
		}

	}

	a.logger.WithField("operation", operation.Operation).Info("migration complete")

	return nil

}

func main() {

	ctx := context.TODO()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to load aws default config")
	}

	loadConfig(ctx, awsCfg)

	_, db, err := initDB()
	if err != nil {
		logger.WithError(err).Error("failed to initialize database")
	}

	app := app{
		logger: logger,
		db:     db,
	}

	lambda.Start(app.handle)

}
