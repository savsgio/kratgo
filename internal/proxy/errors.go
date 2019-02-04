package proxy

import "errors"

var errMissingLocation = errors.New("missing Location header for http redirect")
var errTooManyRedirects = errors.New("too many redirects detected when doing the request")
