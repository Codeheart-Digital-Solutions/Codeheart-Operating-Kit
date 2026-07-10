package operatingkit

import "embed"

// EmbeddedResources carries source Operating Kit content identity and doctrine into the Go CLI.
// Release catalogs are generated after archives exist and are deliberately not embedded.
//
//go:embed components profiles templates schemas manifest.yaml
var EmbeddedResources embed.FS
