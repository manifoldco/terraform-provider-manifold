package manifold

import "errors"

var (
	errAPITokenRequired = errors.New("An API Token is required to use this provider")
	errProjectNotFound  = errors.New("Could not find project for label")
)
