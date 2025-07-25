package db

import (
	"embed"
)

//go:embed migration/*.sql
var Migration embed.FS
