package frontend

import (
	"errors"
)

// URL is where the CLI makes http requests to
const URL = "http://localhost:7274"

// ErrRequestFailed occurs when failed to make an http request
var ErrRequestFailed = errors.New("cli: failed to make an http request, did you start realmsd?")

// ErrReadResponseFailed occurs when failed to read the http response
var ErrReadResponseFailed = errors.New("cli: failed to read the http response")

// ErrInvalidResponse occurs when the http response is invalid
var ErrInvalidResponse = errors.New("cli: failed to parse the http response")
