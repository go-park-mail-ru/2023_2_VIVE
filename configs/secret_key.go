package configs

import "os"

var (
        SECRET_KEY = os.Getenv("SECRET_KEY")
)
