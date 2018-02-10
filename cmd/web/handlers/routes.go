package handlers

import (
	"net/http"
	"revelforce/cmd/web/middleware"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func Routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", dashboard).Methods("GET").Name("admin.dashboard")

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/recover", resetPasswordForm).Queries("email", "{email}").Queries("hash", "{hash}").Methods("GET")
	auth.HandleFunc("/reset", postPasswordReset).Methods("POST")
	auth.HandleFunc("/forgot", forgotPasswordForm).Methods("GET")
	auth.HandleFunc("/forgot", postForgotPassword).Methods("POST")
	auth.HandleFunc("/signup", signupForm).Methods("GET")
	auth.HandleFunc("/signup", postSignup).Methods("POST")
	auth.HandleFunc("/login", loginForm).Methods("GET")
	auth.HandleFunc("/login", postLogin).Methods("POST")

	user := r.PathPrefix("/u").Subrouter()
	user.HandleFunc("/profile", profileForm).Methods("GET")
	user.HandleFunc("/profile", postProfile).Methods("POST")
	user.HandleFunc("/password", passwordForm).Methods("GET")
	user.HandleFunc("/password", postPassword).Methods("POST")
	user.HandleFunc("/logout", logout).Methods("GET")

	admin.HandleFunc("/trip", removeTrip).Queries("remove", "").Methods("GET").Name("admin.trip.delete")
	admin.HandleFunc("/trip", tripForm).Methods("GET").Name("admin.trip")
	admin.HandleFunc("/trip", postTrip).Methods("POST").Name("admin.trip")
	admin.HandleFunc("/trips", listTrips).Methods("GET").Name("admin.trips")

	admin.HandleFunc("/vendor", removeVendor).Queries("remove", "").Methods("GET").Name("admin.vendor.delete")
	admin.HandleFunc("/vendor", vendorForm).Methods("GET").Name("admin.vendor")
	admin.HandleFunc("/vendor", postVendor).Methods("POST").Name("admin.vendor")
	admin.HandleFunc("/vendors", listVendors).Methods("GET").Name("admin.vendors")

	admin.HandleFunc("/user", removeUser).Queries("remove", "").Methods("GET").Name("admin.user.delete")
	admin.HandleFunc("/user", userForm).Methods("GET").Name("admin.user")
	admin.HandleFunc("/user", postUser).Methods("POST").Name("admin.user")
	admin.HandleFunc("/users", listUsers).Methods("GET").Name("admin.users")

	fs := http.FileServer(http.Dir(viper.GetString("files.static")))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.NotFoundHandler = http.HandlerFunc(notFound)

	sirMuxalot := http.NewServeMux()
	sirMuxalot.Handle("/", r)

	sirMuxalot.Handle("/u/", negroni.New(
		negroni.HandlerFunc(requireLogin),
		negroni.Wrap(r),
	))

	sirMuxalot.Handle("/admin/", negroni.New(
		negroni.HandlerFunc(requireAdmin),
		negroni.Wrap(r),
	))

	sirMuxalot.Handle("/auth/", negroni.New(
		negroni.HandlerFunc(requireGuest),
		negroni.Wrap(r),
	))

	n := negroni.New()
	n.UseHandler(sirMuxalot)
	return middleware.SecureHeaders(middleware.NoSurf(n))
}
