package view

import (
	"log"
	"net/http"
	"runtime/debug"
)

func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	Render(w, r, "error", &View{
		Err: appError{
			Code:    status,
			Message: http.StatusText(status),
		},
	})
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	//stack trace appended to logging of error
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	Render(w, r, "error", &View{
		Err: appError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		},
	})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Render(w, r, "error", &View{
		Err: appError{
			Code:    404,
			Message: "Not found",
		},
	})
}
