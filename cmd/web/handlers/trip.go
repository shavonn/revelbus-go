package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/flash"
	"strconv"

	"github.com/gorilla/mux"
)

func TripForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "admin-trip", &view.View{
			Form:  new(models.TripForm),
			Title: "New Trip",
		})
		return
	}

	t := &models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.Get()
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.GetVendors(true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.TripForm{
		ID:           id,
		Title:        t.Title,
		Slug:         t.Slug,
		Status:       t.Status,
		Blurb:        t.Blurb,
		Description:  t.Description,
		Start:        t.Start.Format(db.TimeFormat),
		End:          t.End.Format(db.TimeFormat),
		Price:        t.Price,
		TicketingURL: t.TicketingURL,
		Notes:        t.Notes,
		Image:        t.Image,
	}

	view.Render(w, r, "admin-trip", &view.View{
		ActiveKey: "trip",
		Form:      f,
		Trip:      t,
		Vendors:   vendors,
	})
}

func PostTrip(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.TripForm{
		ID:           r.PostForm.Get("id"),
		Title:        r.PostForm.Get("title"),
		Slug:         r.PostForm.Get("slug"),
		Status:       r.PostForm.Get("status"),
		Blurb:        r.PostForm.Get("blurb"),
		Description:  r.PostForm.Get("description"),
		Start:        r.PostForm.Get("start"),
		End:          r.PostForm.Get("end"),
		TicketingURL: r.PostForm.Get("ticketing_url"),
		Price:        r.PostForm.Get("price"),
		Notes:        r.PostForm.Get("notes"),
		Image:        r.PostForm.Get("image"),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if id == "" {
			v.Title = "New Trip"
		}

		view.Render(w, r, "admin-trip", v)
	}

	fn, err := utils.UploadFile(w, r, "trip_image", "uploads/trip/")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(fn) > 0 && fn[0] != "" {
		f.Image = fn[0]
	} else if (len(f.Image) != 0) && (len(r.Form["deleteimg"]) == 1) {
		err = utils.DeleteFile(f.Image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		f.Image = ""
	}

	var msg string

	t := models.Trip{
		Title:        f.Title,
		Slug:         f.Slug,
		Status:       f.Status,
		Blurb:        f.Blurb,
		Description:  f.Description,
		Start:        db.ToTime(f.Start),
		End:          db.ToTime(f.End),
		TicketingURL: f.TicketingURL,
		Price:        f.Price,
		Notes:        f.Notes,
		Image:        f.Image,
	}

	if id != "" {
		t.ID = utils.ToInt(id)

		err := t.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		msg = utils.MsgSuccessfullyUpdated
	} else {
		err := t.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		id = strconv.Itoa(t.ID)
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trip?id="+id, http.StatusSeeOther)
}

func ListTrips(w http.ResponseWriter, r *http.Request) {
	trips, err := models.GetTrips()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "admin-trips", &view.View{
		Title: "Trips",
		Trips: trips,
	})
}

func RemoveTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t := &models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trips", http.StatusSeeOther)
}

func TripPartners(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t := &models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.Get()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.GetVendors(true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "trip-partners", &view.View{
		ActiveKey: "partners",
		Trip:      t,
		Vendors:   vendors,
	})
}

func TripVenues(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t := &models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.Get()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.GetVendors(true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "trip-venues", &view.View{
		ActiveKey: "venues",
		Trip:      t,
		Vendors:   vendors,
	})
}

func AttachVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	vid := r.PostForm.Get("vendor")
	role := r.PostForm.Get("role")

	t := models.Trip{
		ID: utils.ToInt(id),
	}

	err = t.AttachVendor(role, vid)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyAddedVendor, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trip/"+id+"?"+role+"s", http.StatusSeeOther)
}

func DetachVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vid := vars["vid"]
	role := vars["role"]

	t := models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.DetachVendor(role, vid)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemovedVendor, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trip/"+id+"?"+role+"s", http.StatusSeeOther)
}

func UpdateVenueStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vid := vars["vid"]
	isPrimary, _ := strconv.ParseBool(vars["is_primary"])

	t := models.Trip{
		ID: utils.ToInt(id),
	}

	err := t.SetVenueStatus(vid, isPrimary)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyUpdated, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/trip/"+id+"?venues", http.StatusSeeOther)
}
