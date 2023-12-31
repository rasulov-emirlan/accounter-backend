package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	server struct {
		Port           string   `env:"PORT" env-default:":8080"`
		AllowedOrigins []string `env:"ALLOWED_CORS_ORIGINS" env-default:"*"`
		// TODO: implement reader
		// that would read these two
		// fields from a yaml file
		TimeoutRead  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"15s"`
		TimeoutWrite time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"15s"`
	}

	flags struct {
		envFilename    string
		DevMode        bool
		WithMigrations bool
	}

	Config struct {
		Server      server
		LogLevel    string `env:"LOG_LEVEL" env-default:"debug"`
		ServiceName string `env:"SERVICE_NAME" env-default:"accounter-backend"`
		JWTsecret   string `env:"JWT_SECRET" env-default:"supersecret"`
		DatabaseURL string `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@localhost:5432/accounter?sslmode=disable"`
		JeagerURL   string `env:"JAEGER_URL" env-default:"http://localhost:14268/api/traces"`
		Flags       flags

		Usage func()
	}
)

func LoadConfig() (Config, error) {
	var cfg Config

	cfg.Flags = loadFlags()
	if cfg.Flags.DevMode {
		cfg.LogLevel = "dev"
	}

	header := "Config loaded from file"

	cfg.Usage = cleanenv.FUsage(os.Stdout, cfg, &header)

	if cfg.Flags.envFilename != "" {
		if err := cleanenv.ReadConfig(cfg.Flags.envFilename, &cfg); err != nil {
			return Config{}, err
		}
		port(&cfg)
		return cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.Flags.DevMode {
		cfg.LogLevel = "dev"
	}

	port(&cfg)
	return cfg, nil
}

func port(cfg *Config) {
	// remove : from port
	cfg.Server.Port = strings.TrimPrefix(cfg.Server.Port, ":")
	cfg.Server.Port = fmt.Sprintf(":%s", cfg.Server.Port)
}

func loadFlags() flags {
	var f flags

	flag.BoolVar(&f.DevMode, "dev", false, "Run in dev mode, some features will be disabled.\nFor example, emails will be printed to stdout instead of being sent.")
	flag.BoolVar(&f.WithMigrations, "migrate", false, "Start application with db migrations.")
	flag.StringVar(&f.envFilename, "env", "", "Path to .env file.")
	flag.Parse()

	return f
}
