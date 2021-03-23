package utils

import (
	"net"
	"net/http"
	"reflect"

	"github.com/Azure/go-autorest/autorest"
)

func ResponseWasNotFound(resp autorest.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

// err here will be a "github.com/Azure/azure-sdk-for-go/sdk/internal/runtime".ResponseError,
// but it is in an internal package, we cannot directly use it,
// therefore we have to use reflect to get the raw http response out of it
func Track2ResponseWasNotFound(err interface{}) bool {
	resp := reflect.ValueOf(err).MethodByName("RawResponse").Call([]reflect.Value{})[0].Interface().(*http.Response)
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
