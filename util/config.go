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
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Try to read the config file
	err = viper.ReadInConfig()
	if err != nil {
		// Check if the error is specifically about the config file not being found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, but that's okay - we'll use environment variables
			// Log this if you want to know when fallback is happening
			// fmt.Println("Config file not found, using environment variables only")
		} else {
			// Config file was found but another error occurred (e.g., parsing error)
			return config, err
		}
	}

	// Unmarshal the configuration (from file if found, or from env vars only)
	err = viper.Unmarshal(&config)
	return config, err
}
