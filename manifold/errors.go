package manifold

import "errors"

var (
	errAPITokenRequired = errors.New("an API Token is required to use this provider")
	errProjectNotFound  = errors.New("could not find project for label")
)
