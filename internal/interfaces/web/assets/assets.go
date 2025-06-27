// Package assets provides the assets for the web frontend.
package assets

import (
	"embed"
)

// ABOUT: This embeds the files directly into the binary during compilation.
// If you need to add any scripts you can do so by adding them to the scripts directory.
// Remember to uncomment the embed directive below when adding new files.

// //go:embed scripts/*

// Content holds the embedded assets.
var Content embed.FS
