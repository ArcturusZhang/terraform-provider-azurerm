package utils

import (
	"net"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

func ResponseWasNotFound(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

func Track2ResponseWasNotFound(resp *http.Response) bool {
	return HTTPResponseWasStatusCode(resp, http.StatusNotFound)
}

func ResponseWasBadRequest(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusBadRequest)
}

func ResponseWasForbidden(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusForbidden)
}

func ResponseWasConflict(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusConflict)
}

func ResponseErrorIsRetryable(err error) bool {
	if arerr, ok := err.(autorest.DetailedError); ok {
		err = arerr.Original
	}

	// nolint gocritic
	switch e := err.(type) {
	case net.Error:
		if e.Temporary() || e.Timeout() {
			return true
		}
	}

	return false
}

func ResponseWasStatusCode(resp autorest.Response, statusCode int) bool { // nolint: unparam
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}

func HTTPResponseWasStatusCode(resp *http.Response, statusCode int) bool { // nolint: unparam
	if resp != nil {
		if resp.StatusCode == statusCode {
			return true
		}
	}

	return false
}
