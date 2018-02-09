package handlers

import (
	"net/http"
	"revelforce-admin/handlers/middleware"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func Routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/tour", tourForm).Methods("GET").Name("admin.tour")

	fs := http.FileServer(http.Dir(viper.GetString("files.static")))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	sirMuxalot := http.NewServeMux()
	sirMuxalot.Handle("/", r)

	n := negroni.New()
	n.UseHandler(sirMuxalot)
	return middleware.SecureHeaders(middleware.NoSurf(n))
}
