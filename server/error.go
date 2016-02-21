package server

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := f(w, r, p)
		if err == nil {
			return
		}

		switch err.(type) {
		case badRequest:
			handleOutput(w, http.StatusBadRequest, err.Error())
		case notFound:
			handleOutput(w, http.StatusNotFound, "not found")
		default:
			log.Println(err)
			handleOutput(w, http.StatusInternalServerError, "oops")
		}
	}
}
