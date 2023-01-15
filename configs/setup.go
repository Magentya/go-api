package configs

import "github.com/joho/godotenv"

func Setup() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	return nil
}
