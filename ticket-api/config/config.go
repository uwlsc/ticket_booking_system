package config

import (
	"os"

	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

type Config struct {
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`

	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBType     string `mapstructure:"DB_TYPE"`

	MailClientID     string `mapstructure:"MAIL_CLIENT_ID"`
	MailClientSecret string `mapstructure:"MAIL_CLIENT_SECRET"`
	MailTokenType    string `mapstructure:"MAIL_TOKEN_TYPE"`

	SentryDSN          string `mapstructure:"SENTRY_DSN"`
	MaxMultipartMemory int64  `mapstructure:"MAX_MULTIPART_MEMORY"`
	StorageBucketName  string `mapstructure:"STORAGE_BUCKET_NAME"`
	StripeKey          string `mapstructure:"STRIPE_KEY"`

	TimeZone              string `mapstructure:"TIMEZONE"`
	AdminEmail            string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword         string `mapstructure:"ADMIN_PASSWORD"`
	RecaptchaSecretKey    string `mapstructure:"RECAPTCHA_SECRET_KEY"`
	GeoCodingAPIKey       string `mapstructure:"GEO_CODING_API_KEY"`
	GCPProjectID          string `mapstructure:"GCP_PROJECT_ID"`
	STRIPE_WEBHOOK_SECRET string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}

var globalEnv = Config{
	MaxMultipartMemory: 10 << 20, // 10 MB
}

func GetEnv() Config {
	return globalEnv
}

var jwtConfig *jwt.Config

func GetJWTConfig() *jwt.Config {
	return jwtConfig
}

func New() Config {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}

	jsonKey, err := os.ReadFile("./config/serviceAccountKey.json")
	if err != nil {
		os.Exit(1)
	}

	jwtConfig, err = google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		os.Exit(1)
	}

	viper.SetDefault("TIMEZONE", "UTC")

	err = viper.Unmarshal(&globalEnv)
	if err != nil {
		os.Exit(1)
	}

	return globalEnv
}
