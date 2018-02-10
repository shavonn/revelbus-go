package handlers

import (
	"net/http"
	"revelforce/cmd/web/handlers/middleware"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func Routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", dashboard).Methods("GET").Name("admin.dashboard")

	admin.HandleFunc("/trip", removeTrip).Queries("remove", "").Methods("GET").Name("admin.trip.delete")
	admin.HandleFunc("/trip", tripForm).Methods("GET").Name("admin.trip")
	admin.HandleFunc("/trip", postTrip).Methods("POST").Name("admin.trip")
	admin.HandleFunc("/trips", listTrips).Methods("GET").Name("admin.trips")

	admin.HandleFunc("/vendor", removeVendor).Queries("remove", "").Methods("GET").Name("admin.vendor.delete")
	admin.HandleFunc("/vendor", vendorForm).Methods("GET").Name("admin.vendor")
	admin.HandleFunc("/vendor", postVendor).Methods("POST").Name("admin.vendor")
	admin.HandleFunc("/vendors", listVendors).Methods("GET").Name("admin.vendors")

	fs := http.FileServer(http.Dir(viper.GetString("files.static")))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	sirMuxalot := http.NewServeMux()
	sirMuxalot.Handle("/", r)

	n := negroni.New()
	n.UseHandler(sirMuxalot)
	return middleware.SecureHeaders(middleware.NoSurf(n))
}
