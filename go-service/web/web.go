package web

import (
	"embed"
)

//go:embed static
var Static embed.FS

//go:embed pages
var Pages embed.FS
