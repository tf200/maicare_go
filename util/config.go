package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbSource              string        `mapstructure:"DB_SOURCE"`
	ServerAddress         string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenSecretKey  string        `mapstructure:"ACCESS_TOKEN_SECRET_KEY"`
	AccessTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenSecretKey string        `mapstructure:"REFRESH_TOKEN_SECRET_KEY"`
	RefreshTokenDuration  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TwoFATokenSecretKey   string        `mapstructure:"2FA_TOKEN_SECRET_KEY"`
	TwoFATokenDuration    time.Duration `mapstructure:"2FA_TOKEN_DURATION"`
	B2Key                 string        `mapstructure:"B2_KEY"`
	B2KeyID               string        `mapstructure:"B2_KEY_ID"`
	B2Bucket              string        `mapstructure:"B2_BUCKET"`
	Host                  string        `mapstructure:"HOST"`
	RedisHost             string        `mapstructure:"REDIS_HOST"`
	RedisPassword         string        `mapstructure:"REDIS_PASSWORD"`
	Remote                bool          `mapstructure:"REMOTE"`
	OpenRouterAPIKey      string        `mapstructure:"OPEN_ROUTER_API_KEY"`
	SmtpName              string        `mapstructure:"SMTP_NAME"`
	SmtpAddress           string        `mapstructure:"SMTP_ADDRESS"`
	SmtpAuth              string        `mapstructure:"SMTP_AUTH"`
	SmtpHost              string        `mapstructure:"SMTP_HOST"`
	SmtpPort              int           `mapstructure:"SMTP_PORT"`
	BrevoSenderName       string        `mapstructure:"BREVO_SENDER_NAME"`
	BrevoSenderEmail      string        `mapstructure:"BREVO_SENDER_EMAIL"`
	BrevoApiKey           string        `mapstructure:"BREVO_API_KEY"`
	Environment           string        `mapstructure:"ENVIRONMENT"`
	GrpcUrl               string        `mapstructure:"GRPC_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Bind all environment variables explicitly
	// This works for both file-based config and environment variables
	envVars := []string{
		"DB_SOURCE", "SERVER_ADDRESS", "ACCESS_TOKEN_SECRET_KEY",
		"ACCESS_TOKEN_DURATION", "REFRESH_TOKEN_SECRET_KEY",
		"REFRESH_TOKEN_DURATION", "2FA_TOKEN_SECRET_KEY",
		"2FA_TOKEN_DURATION", "B2_KEY", "B2_KEY_ID", "B2_BUCKET",
		"HOST", "REDIS_HOST", "REDIS_PASSWORD", "REMOTE",
		"OPEN_ROUTER_API_KEY", "SMTP_NAME", "SMTP_ADDRESS",
		"SMTP_AUTH", "SMTP_HOST", "SMTP_PORT", "BREVO_SENDER_NAME",
		"BREVO_SENDER_EMAIL", "BREVO_API_KEY",
	}

	for _, envVar := range envVars {
		viper.BindEnv(envVar)
	}

	// Try to read the config file (works locally)
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found (normal in Docker), using environment variables only
		} else {
			// Config file was found but another error occurred
			return config, err
		}
	}

	// Unmarshal the configuration
	err = viper.Unmarshal(&config)
	return config, err
}
