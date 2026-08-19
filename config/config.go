package config

import (
	"fmt"
	"strings"
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
	TwoFATokenSecretKey   string        `mapstructure:"TWO_FA_TOKEN_SECRET_KEY"`
	TwoFATokenDuration    time.Duration `mapstructure:"TWO_FA_TOKEN_DURATION"`
	B2Endpoint            string        `mapstructure:"B2_ENDPOINT"`
	B2Key                 string        `mapstructure:"B2_KEY"`
	B2KeyID               string        `mapstructure:"B2_KEY_ID"`
	B2Bucket              string        `mapstructure:"B2_BUCKET"`
	DisableBucket         bool          `mapstructure:"DISABLE_BUCKET"`
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
	GrpcURL               string        `mapstructure:"GRPC_URL"`
	OpenRouterApiKey      string        `mapstructure:"OPEN_ROUTER_API_KEY"`
	OpenRouterModel       string        `mapstructure:"OPEN_ROUTER_MODEL"`
	MigrationsPath        string        `mapstructure:"MIGRATIONS_PATH"`
	AdminEmail            string        `mapstructure:"ADMIN_EMAIL"`
	AdminPassword         string        `mapstructure:"ADMIN_PASSWORD"`
	SystemActorUserID     string        `mapstructure:"SYSTEM_ACTOR_USER_ID"`
	SystemActorEmployeeID string        `mapstructure:"SYSTEM_ACTOR_EMPLOYEE_ID"`
	WsAllowedOrigins      string        `mapstructure:"WS_ALLOWED_ORIGINS"`
	WsTicketTTL           time.Duration `mapstructure:"WS_TICKET_TTL"`
}

func Load(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetDefault("WS_TICKET_TTL", "1m")
	viper.AutomaticEnv()

	envVars := []string{
		"DB_SOURCE", "SERVER_ADDRESS", "ACCESS_TOKEN_SECRET_KEY",
		"ACCESS_TOKEN_DURATION", "REFRESH_TOKEN_SECRET_KEY",
		"REFRESH_TOKEN_DURATION", "TWO_FA_TOKEN_SECRET_KEY",
		"TWO_FA_TOKEN_DURATION", "B2_ENDPOINT", "B2_KEY", "B2_KEY_ID", "B2_BUCKET",
		"DISABLE_BUCKET", "HOST", "REDIS_HOST", "REDIS_PASSWORD", "REMOTE",
		"OPEN_ROUTER_API_KEY", "SMTP_NAME", "SMTP_ADDRESS",
		"SMTP_AUTH", "SMTP_HOST", "SMTP_PORT", "BREVO_SENDER_NAME",
		"BREVO_SENDER_EMAIL", "BREVO_API_KEY", "ENVIRONMENT", "GRPC_URL",
		"OPEN_ROUTER_API_KEY", "OPEN_ROUTER_MODEL",
		"MIGRATIONS_PATH", "ADMIN_EMAIL", "ADMIN_PASSWORD",
		"SYSTEM_ACTOR_USER_ID", "SYSTEM_ACTOR_EMPLOYEE_ID",
		"WS_ALLOWED_ORIGINS", "WS_TICKET_TTL",
	}

	for _, envVar := range envVars {
		err = viper.BindEnv(envVar)
		if err != nil {
			return cfg, fmt.Errorf("failed to bind env var %s: %w", envVar, err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, err
		}
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return cfg, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	crucialVars := map[string]string{
		"DB_SOURCE":                cfg.DbSource,
		"SERVER_ADDRESS":           cfg.ServerAddress,
		"ACCESS_TOKEN_SECRET_KEY":  cfg.AccessTokenSecretKey,
		"REFRESH_TOKEN_SECRET_KEY": cfg.RefreshTokenSecretKey,
		"TWO_FA_TOKEN_SECRET_KEY":  cfg.TwoFATokenSecretKey,
		"TWO_FA_TOKEN_DURATION":    cfg.TwoFATokenDuration.String(),
		"HOST":                     cfg.Host,
		"ENVIRONMENT":              cfg.Environment,
		"GRPC_URL":                 cfg.GrpcURL,
		"MIGRATIONS_PATH":          cfg.MigrationsPath,
		"SYSTEM_ACTOR_USER_ID":     cfg.SystemActorUserID,
		"SYSTEM_ACTOR_EMPLOYEE_ID": cfg.SystemActorEmployeeID,
	}

	var missingVars []string
	for varName, varValue := range crucialVars {
		if strings.TrimSpace(varValue) == "" {
			missingVars = append(missingVars, varName)
		}
	}

	if cfg.AccessTokenDuration <= 0 {
		missingVars = append(missingVars, "ACCESS_TOKEN_DURATION")
	}
	if cfg.RefreshTokenDuration <= 0 {
		missingVars = append(missingVars, "REFRESH_TOKEN_DURATION")
	}
	if cfg.TwoFATokenDuration <= 0 {
		missingVars = append(missingVars, "TWO_FA_TOKEN_DURATION")
	}
	if cfg.WsTicketTTL <= 0 {
		missingVars = append(missingVars, "WS_TICKET_TTL")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing or invalid crucial environment variables: %s", strings.Join(missingVars, ", "))
	}

	return nil
}
