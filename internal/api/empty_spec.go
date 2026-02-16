//go:build disable_spec

package api

// Spec returns nil when built with -tags=disable_spec,
// disabling the Swagger UI and spec serving middleware.
func Spec() []byte {
	return nil
}
