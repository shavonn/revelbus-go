package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strconv"
)

func tripForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		render(w, r, "trip.html", &view{
			Form:  new(forms.TripForm),
			Title: "New Trip",
		})
		return
	}

	t, err := db.GetTripByID(id)
	if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.TripForm{
		ID:           id,
		Title:        t.Title,
		Slug:         t.Slug,
		Status:       t.Status,
		Description:  t.Description,
		Start:        t.Start.Format(db.TimeFormat),
		End:          t.End.Format(db.TimeFormat),
		TicketingURL: t.TicketingURL,
		Notes:        t.Notes,
		Image:        t.Image,
	}

	render(w, r, "trip.html", &view{
		Form: f,
		Trip: *t,
	})
}

func postTrip(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.TripForm{
		ID:           r.PostForm.Get("id"),
		Title:        r.PostForm.Get("title"),
		Slug:         r.PostForm.Get("slug"),
		Status:       r.PostForm.Get("status"),
		Description:  r.PostForm.Get("description"),
		Start:        r.PostForm.Get("start"),
		End:          r.PostForm.Get("end"),
		TicketingURL: r.PostForm.Get("ticketing_url"),
		Notes:        r.PostForm.Get("notes"),
		Image:        r.PostForm.Get("image"),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		if id == "" {
			v.Title = "New Trip"
		}

		render(w, r, "trip.html", v)
	}

	fn, err := uploadFile(w, r, "trip_image", "uploads/trip/")
	if err != nil {
		serverError(w, r, err)
		return
	}

	if fn != "" {
		f.Image = fn
	} else if (len(f.Image) != 0) && (len(r.Form["deleteimg"]) == 1) {
		err = deleteFile("uploads/trip/" + f.Image)
		if err != nil {
			serverError(w, r, err)
			return
		}
		f.Image = ""
	}

	var msg string

	trip := db.Trip{
		Title:        f.Title,
		Slug:         f.Slug,
		Status:       f.Status,
		Description:  f.Description,
		Start:        db.ToTime(f.Start),
		End:          db.ToTime(f.End),
		TicketingURL: f.TicketingURL,
		Notes:        f.Notes,
		Image:        f.Image,
	}

	if id != "" {
		intID, _ := strconv.Atoi(id)
		trip.ID = intID

		err := trip.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		msg = MsgSuccessfullyUpdated
	} else {
		tid, err := trip.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		id = strconv.Itoa(tid)
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trip?id="+id, http.StatusSeeOther)
}

func listTrips(w http.ResponseWriter, r *http.Request) {
	trips, err := db.GetTrips()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "trips.html", &view{
		Title: "Trips",
		Trips: trips,
	})
}

func removeTrip(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	t, err := db.GetTripByID(id)
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = t.Delete()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyRemoved, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trips", http.StatusSeeOther)
}
