// Package view embeds the tactical front-end assets (CSS/JS) so the Go binary
// serves them regardless of the container working directory.
package view

import "embed"

//go:embed assets/*
var Assets embed.FS
