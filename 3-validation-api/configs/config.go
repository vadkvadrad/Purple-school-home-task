package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Db dbConfig
}

type dbConfig struct {
	Dsn string
}

func Load() *Config {
	err := godotenv.Load(dir(".env"))
	if err != nil {
		log.Println("Error loading .env file, using default config.", "Error:", err.Error())
	}

	return &Config{
		Db: dbConfig{
			Dsn: os.Getenv("DSN"),
		},
	}

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