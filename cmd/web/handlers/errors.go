package handlers

import (
	"log"
	"net/http"
	"runtime/debug"
)

func clientError(w http.ResponseWriter, r *http.Request, status int) {
	render(w, r, "error", &view{
		Err: appError{
			Code:    status,
			Message: http.StatusText(status),
		},
	})
}

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	//stack trace appended to logging of error
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	render(w, r, "error", &view{
		Err: appError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		},
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, "error", &view{
		Err: appError{
			Code:    404,
			Message: "Not found",
		},
	})
}
