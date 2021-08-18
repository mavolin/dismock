// Package assets provides access to the assets of the application.
package assets

import (
	"embed"
)

//go:embed templates
var FS embed.FS
