package web

import (
	"net/http"
	"revelforce/cmd/web/handlers"
	"revelforce/cmd/web/middleware"
	"revelforce/cmd/web/view"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func Routes() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", handlers.Index).Methods("GET")
	r.HandleFunc("/trips", handlers.Trips).Methods("GET")
	r.HandleFunc("/trip/{slug}", handlers.Trip).Methods("GET")
	r.HandleFunc("/faq", handlers.Faq).Methods("GET")
	r.HandleFunc("/about", handlers.About).Methods("GET")
	r.HandleFunc("/contact", handlers.Contact).Methods("GET")
	r.HandleFunc("/contact", handlers.ContactPost).Methods("POST")

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/recover", handlers.ResetPasswordForm).Queries("email", "{email}").Queries("hash", "{hash}").Methods("GET")
	auth.HandleFunc("/reset", handlers.PostPasswordReset).Methods("POST")
	auth.HandleFunc("/forgot", handlers.ForgotPasswordForm).Methods("GET")
	auth.HandleFunc("/forgot", handlers.PostForgotPassword).Methods("POST")
	auth.HandleFunc("/signup", handlers.SignupForm).Methods("GET")
	auth.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	auth.HandleFunc("/login", handlers.LoginForm).Methods("GET")
	auth.HandleFunc("/login", handlers.PostLogin).Methods("POST")

	user := r.PathPrefix("/u").Subrouter()
	user.HandleFunc("/", handlers.UserDashboard).Methods("GET")
	user.HandleFunc("/profile", handlers.ProfileForm).Methods("GET")
	user.HandleFunc("/profile", handlers.PostProfile).Methods("POST")
	user.HandleFunc("/password", handlers.PasswordForm).Methods("GET")
	user.HandleFunc("/password", handlers.PostPassword).Methods("POST")
	user.HandleFunc("/logout", handlers.Logout).Methods("GET")

	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/", handlers.AdminDashboard).Methods("GET")

	admin.HandleFunc("/settings", handlers.SettingsForm).Methods("GET")
	admin.HandleFunc("/settings", handlers.PostSettings).Methods("POST")

	admin.HandleFunc("/file/{id}", handlers.RemoveFile).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/files", handlers.ListFiles).Methods("GET")
	admin.HandleFunc("/upload", handlers.UploadForm).Methods("GET")
	admin.HandleFunc("/upload", handlers.PostUpload).Methods("POST")

	// trip funcs
	admin.HandleFunc("/trip/{id}", handlers.UpdateVenueStatus).Queries("venue", "{vid}").Queries("is_primary", "{is_primary}").Methods("GET")
	admin.HandleFunc("/trip/{id}", handlers.AttachVendor).Queries("vendor", "").Methods("POST")
	admin.HandleFunc("/trip/{id}", handlers.DetachVendor).Queries("vendor", "{vid}").Queries("role", "{role}").Methods("GET")
	admin.HandleFunc("/trip/{id}", handlers.TripVenues).Queries("venues", "").Methods("GET")
	admin.HandleFunc("/trip/{id}", handlers.TripPartners).Queries("partners", "").Methods("GET")

	// trip crud
	admin.HandleFunc("/trip/{id}", handlers.RemoveTrip).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/trip", handlers.TripForm).Methods("GET")
	admin.HandleFunc("/trip", handlers.PostTrip).Methods("POST")
	admin.HandleFunc("/trips", handlers.ListTrips).Methods("GET")

	// vendor crud
	admin.HandleFunc("/vendor/{id}", handlers.RemoveVendor).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/vendor", handlers.VendorForm).Methods("GET")
	admin.HandleFunc("/vendor", handlers.PostVendor).Methods("POST")
	admin.HandleFunc("/vendors", handlers.ListVendors).Methods("GET")

	// faq crud
	admin.HandleFunc("/faq/{id}", handlers.RemoveFAQ).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/faq", handlers.FaqForm).Methods("GET")
	admin.HandleFunc("/faq", handlers.PostFAQ).Methods("POST")
	admin.HandleFunc("/faqs", handlers.ListFAQs).Methods("GET")

	// slide crud
	admin.HandleFunc("/slide/{id}", handlers.RemoveSlide).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/slide", handlers.SlideForm).Methods("GET")
	admin.HandleFunc("/slide", handlers.PostSlide).Methods("POST")
	admin.HandleFunc("/slides", handlers.ListSlides).Methods("GET")

	//user crud
	admin.HandleFunc("/user/{id}", handlers.RemoveUser).Queries("remove", "").Methods("GET")
	admin.HandleFunc("/user", handlers.UserForm).Methods("GET")
	admin.HandleFunc("/user", handlers.PostUser).Methods("POST")
	admin.HandleFunc("/users", handlers.ListUsers).Methods("GET")

	fs := http.FileServer(http.Dir(viper.GetString("files.static")))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.NotFoundHandler = http.HandlerFunc(view.NotFound)

	sirMuxalot := http.NewServeMux()
	sirMuxalot.Handle("/", r)

	sirMuxalot.Handle("/u/", negroni.New(
		negroni.HandlerFunc(middleware.RequireLogin),
		negroni.Wrap(r),
	))

	sirMuxalot.Handle("/admin/", negroni.New(
		negroni.HandlerFunc(middleware.RequireAdmin),
		negroni.Wrap(r),
	))

	sirMuxalot.Handle("/auth/", negroni.New(
		negroni.HandlerFunc(middleware.RequireGuest),
		negroni.Wrap(r),
	))

	n := negroni.New()
	n.UseHandler(sirMuxalot)
	return middleware.SecureHeaders(middleware.NoSurf(n))
}
