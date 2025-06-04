package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Servers     ServerConfig     `mapstructure:"servers"`
	Database    DatabaseConfig   `mapstructure:"database"`
	Stripe      StripeConfig     `mapstructure:"stripe"`
	Environment string           `mapstructure:"environment"`
}

type ServerConfig struct {
	GRPC    GRPCServerConfig    `mapstructure:"grpc"`
	REST    RESTServerConfig    `mapstructure:"rest"`
	Webhook WebhookServerConfig `mapstructure:"webhook"`
}

type GRPCServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type RESTServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type WebhookServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Path string `mapstructure:"path"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type StripeConfig struct {
	SecretKey     string `mapstructure:"secret_key"`
	WebhookSecret string `mapstructure:"webhook_secret"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// Load config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Override with environment variables
	viper.AutomaticEnv()

	// Override Stripe settings with environment variables
	if stripeKey := os.Getenv("STRIPE_SECRET_KEY"); stripeKey != "" {
		viper.Set("stripe.secret_key", stripeKey)
	}
	if webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET"); webhookSecret != "" {
		viper.Set("stripe.webhook_secret", webhookSecret)
	}

	// Override database URL if provided
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		viper.Set("database.url", dbURL)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}