package ebitenaudio

import "github.com/hoani/sndit"

// Compile-time interface check.
var _ sndit.Context = (*Context)(nil)
