package env

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("error getting current working directory: %v", err)
		return
	}

	curr := wd
	for {
		path := filepath.Join(curr, ".env")
		if _, err = os.Stat(path); err == nil {
			if err = godotenv.Load(path); err != nil {
				log.Printf("error loading .env file: %v", err)
				return
			}
			return
		}

		path = filepath.Join(curr, "go.mod")
		if _, err = os.Stat(path); err == nil {
			log.Print("reached project root(go.mod) but .env is not found")
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
}

func GetEnv(name string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	return ""
}

func GetEnvAsInt(name string, defaultValue int) int {
	if s := os.Getenv(name); s != "" {
		res, err := strconv.Atoi(s)
		if err != nil {
			log.Printf("error getting env int value: %v", err)
			return 0
		}
		return res
	} else {
		return defaultValue
	}
}
