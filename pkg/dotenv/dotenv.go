// dotenv package is a convenience package that attempts to load .env into process's environment.
package dotenv

import (
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err == nil {
		slog.Info(".env file detected.")
	} else if !errors.Is(err, os.ErrNotExist) {
		slog.Warn(".env file failed: ", err)
	}
}
