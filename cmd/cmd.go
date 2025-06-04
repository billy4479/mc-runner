package cmd

import (
	"fmt"

	"github.com/billy4479/mc-runner/internal"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func Run() error {
	err := godotenv.Load()
	if err != nil {
		err = fmt.Errorf("load .env: %w", err)
		fmt.Printf("warn: %v, proceeding without.\n", err)
	}

	config := internal.Config{}
	err = envconfig.Process("", &config)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	return internal.Run(&config)
}
