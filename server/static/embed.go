package static

import "embed"

//go:embed dist/*
var StaticFS embed.FS
