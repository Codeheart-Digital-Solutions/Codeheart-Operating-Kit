package operatingkit

import "embed"

// EmbeddedResources carries source Operating Kit content into the Go CLI.
//
//go:embed components profiles templates schemas manifest.yaml
var EmbeddedResources embed.FS
