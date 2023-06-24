package config

import (
	"flag"
	"os"
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
		Server          server
		LoggingFilename string `env:"LOGGING_FILENAME" env-default:""`
		LogLevel        string `env:"LOG_LEVEL" env-default:"debug"`
		JWTsecret       string `env:"JWT_SECRET" env-default:"supersecret"`
		DatabaseURL     string `env:"DATABASE_URL" env-default:"postgres://postgres:postgres@localhost:5432/esep?sslmode=disable"`
		JeagerURL       string `env:"JAEGER_URL" env-default:"http://localhost:14268/api/traces"`
		Flags           flags

		Usage func()
	}
)

func LoadConfig() (Config, error) {
	var cfg Config

	cfg.Flags = loadFlags()

	header := "Config loaded from file"

	cfg.Usage = cleanenv.FUsage(os.Stdout, cfg, &header)

	if cfg.Flags.envFilename != "" {
		if err := cleanenv.ReadConfig(cfg.Flags.envFilename, &cfg); err != nil {
			return Config{}, err
		}
		return cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.Flags.DevMode {
		cfg.LogLevel = "dev"
	}

	return cfg, nil
}

func loadFlags() flags {
	var f flags

	flag.BoolVar(&f.DevMode, "dev", false, "Run in dev mode, some features will be disabled.\nFor example, emails will be printed to stdout instead of being sent.")
	flag.BoolVar(&f.WithMigrations, "migrate", false, "Start application with db migrations.")
	flag.StringVar(&f.envFilename, "env", "", "Path to .env file.")
	flag.Parse()

	return f
}
