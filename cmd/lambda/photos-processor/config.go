package main

import (
	"context"
	"photos"

	config "github.com/ddouglas/config-ssm"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
)

var appConfig struct {
	Mode          string `env:"MODE" default:"server"`
	RunMigrations bool   `env:"RUN_MIGRATIONS" default:"false"`
	AppURL        string `env:"APP_URL,required"`
	Auth0         struct {
		CallbackPath string `env:"AUTH0_CALLBACK_PATH,required"`
		ClientID     string `env:"AUTH0_CLIENT_ID,required"`
		ClientSecret string `ssm:"/photos/auth0_client_secret,required"`
		Domain       string `env:"AUTH0_DOMAIN,required"`
	}
	DB struct {
		Host string `env:"DB_HOST,required"`
		Name string `env:"DB_NAME,required"`
		User string `env:"DB_USER,required"`
		Pass string `ssm:"/photos/db_pass,required"`
	}
	Session struct {
		Key string `ssm:"/photos/session_key,required"`
	}
	Buckets struct {
		Photos  string `env:"PHOTOS_BUCKET,required"`
		Uploads string `env:"PHOTOS_UPLOADS_BUCKET,required"`
	}
	Thumbmails struct {
		Width  int `env:"THUMBNAIL_WIDTH,required"`
		Height int `env:"THUMBNAIL_HEIGHT,required"`
	}
	Environment photos.Environment `env:"ENVIRONMENT,required"`
	Server      struct {
		Port string `env:"SERVER_PORT" default:"8080"`
	}
}

func loadConfig(ctx context.Context, awsConfig aws.Config) {

	_ = godotenv.Load()

	ssmClient := ssm.NewFromConfig(awsConfig)

	err := config.Load(
		ctx,
		&appConfig,
		// config.WithPrefix("photos"),
		config.WithSSMClient(ssmClient),
	)
	if err != nil {
		logger.WithError(err).Fatal("failed load configuration")
	}

}
