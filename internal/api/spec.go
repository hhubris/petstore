//go:build !disable_spec

package api

import _ "embed"

//go:embed api.yml
var specYAML []byte

// Spec returns a copy of the embedded OpenAPI spec. Build
// with -tags=disable_spec to exclude the spec from the
// binary (Spec returns nil in that case).
func Spec() []byte {
	out := make([]byte, len(specYAML))
	copy(out, specYAML)
	return out
}
