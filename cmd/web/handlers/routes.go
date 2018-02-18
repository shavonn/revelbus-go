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
	r.HandleFunc("/trips", trips).Methods("GET")
	r.HandleFunc("/trip/{slug}", trip).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/about", about).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/contact", contactPost).Methods("POST")

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
	user.HandleFunc("/", userDashboard).Methods("GET")
	user.HandleFunc("/profile", profileForm).Methods("GET")
	user.HandleFunc("/profile", postProfile).Methods("POST")
	user.HandleFunc("/password", passwordForm).Methods("GET")
	user.HandleFunc("/password", postPassword).Methods("POST")
	user.HandleFunc("/logout", logout).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", adminDashboard).Methods("GET")

	admin.HandleFunc("/settings", settingsForm).Methods("GET")
	admin.HandleFunc("/settings", postSettings).Methods("POST")

	// trip funcs
	admin.HandleFunc("/trip/{id}", updateVenueStatus).Queries("venue", "{vid}").Queries("is_primary", "{is_primary}").Methods("GET")
	admin.HandleFunc("/trip/{id}", attachVendor).Queries("vendor", "").Methods("POST")
	admin.HandleFunc("/trip/{id}", detachVendor).Queries("vendor", "{vid}").Queries("role", "{role}").Methods("GET")
	admin.HandleFunc("/trip/{id}", tripVenues).Queries("venues", "").Methods("GET")
	admin.HandleFunc("/trip/{id}", tripPartners).Queries("partners", "").Methods("GET")

	// trip crud
	admin.HandleFunc("/trip/{id}", removeTrip).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/trip", tripForm).Methods("GET")
	admin.HandleFunc("/trip", postTrip).Methods("POST")
	admin.HandleFunc("/trips", listTrips).Methods("GET")

	// vendor crud
	admin.HandleFunc("/vendor/{id}", removeVendor).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/vendor", vendorForm).Methods("GET")
	admin.HandleFunc("/vendor", postVendor).Methods("POST")
	admin.HandleFunc("/vendors", listVendors).Methods("GET")

	// faq crud
	admin.HandleFunc("/faq/{id}", removeFAQ).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/faq", faqForm).Methods("GET")
	admin.HandleFunc("/faq", postFAQ).Methods("POST")
	admin.HandleFunc("/faqs", listFAQs).Methods("GET")

	// slide crud
	admin.HandleFunc("/slide/{id}", removeSlide).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/slide", slideForm).Methods("GET")
	admin.HandleFunc("/slide", postSlide).Methods("POST")
	admin.HandleFunc("/slides", listSlides).Methods("GET")

	//user crud
	admin.HandleFunc("/user/{id}", removeUser).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/user", userForm).Methods("GET")
	admin.HandleFunc("/user", postUser).Methods("POST")
	admin.HandleFunc("/users", listUsers).Methods("GET")

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
