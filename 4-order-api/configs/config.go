package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Dsn string
}


func Load() (*Config, error) {
	err := godotenv.Load(dir(".env"))
	if err != nil {
		log.Println("Error loading .env file, using default config.", "Error:", err.Error())
	}
	config := &Config{
        Db: DbConfig{
            Dsn: getEnv("DSN", ""),
        },
    }

    if err := config.validate(); err != nil {
        return nil, err
    }

    return config, nil
}

func (c *Config) validate() error {
    if c.Db.Dsn == "" {
        return errors.New("DSN is required")
    }
    return nil
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}