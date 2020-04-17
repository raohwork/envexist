/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

// Package envexist provides few tools to validate env vars based on modules
//
// Each module registers itself to envexist (typically in init()), and waits the
// result in a callback. The application entry (main() in package main) takes
// responsibility to call envexist.Parse(), which validates env vars and triggers
// further initialization.
//
// Package envexist IS NOT (and cannot be) thread-safe.
package envexist // import "github.com/raohwork/envexist"
