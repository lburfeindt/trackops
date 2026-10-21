// Package frontend embeds the static assets served by the desktop app.
package frontend

import "embed"

// Assets contains the static files served by Wails.
//
//go:embed all:dist
var Assets embed.FS
