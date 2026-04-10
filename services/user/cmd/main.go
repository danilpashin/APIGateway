package cmd

import (
	"database/sql"
	"pkg/env"
	"user/internal/config"
)

func main() {
	env.LoadEnv()
}

func runServer(cfg *config.Config, db *sql.DB) {

}
