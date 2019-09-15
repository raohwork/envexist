// Package envexist provides few tools to validate env vars based on modules
//
// Each module registers itself to envexist (typically in init()), and waits the
// result in a callback. The application entry (main() in package main) takes
// responsibility to call envexist.Parse(), which validates env vars and triggers
// further initialization.
//
// Package envexist IS NOT (and cannot be) thread-safe.
package envexist // import "github.com/raohwork/envexist"
