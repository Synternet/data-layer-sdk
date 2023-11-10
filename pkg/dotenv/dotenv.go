// dotenv package is a convenience package that attempts to load .env into process's environment.
package dotenv

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err == nil {
		log.Println(".env file detected.")
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Println(".env file failed: ", err.Error())
	}
}
