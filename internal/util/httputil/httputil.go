package httputil

import (
	"net/http"
)

func BasicAuth(
	next http.Handler,
	validator func(username, password string, r *http.Request) (bool, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //nolint:varnamelen
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.
		username, password, ok := r.BasicAuth()
		if ok {
			if ok, err := validator(username, password, r); err != nil {
				// If an error occurred during validation, then return a
				// 500 error.
				http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: wrap error
			} else if ok {
				// If the username and password are correct, then call
				// the next handler in the chain. Make sure to return
				// afterwards, so that none of the code below is run.
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, `{"message":"Unauthorized"}`, http.StatusUnauthorized)
	}
}
