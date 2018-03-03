package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/domain/models"
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

	err := t.Fetch()
	if err != nil {
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.FetchVendors(true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	galleries, err := models.FetchGalleries()
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
		Start:        t.Start.Format(domain.TimeFormat),
		End:          t.End.Format(domain.TimeFormat),
		Price:        t.Price,
		TicketingURL: t.TicketingURL,
		Notes:        t.Notes,
		ImageID:      t.ImageID,
		GalleryID:    t.GalleryID,
	}

	if t.Image != nil {
		f.Image = t.Image.Thumb.String
	}

	view.Render(w, r, "admin-trip", &view.View{
		ActiveKey: "trip",
		Form:      f,
		Trip:      t,
		Vendors:   vendors,
		Galleries: galleries,
	})
}

func PostTrip(w http.ResponseWriter, r *http.Request) {
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
		ImageID:      utils.ToInt(r.PostForm.Get("image_id")),
		GalleryID:    utils.ToInt(r.PostForm.Get("gallery_id")),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if f.ID == "" {
			v.Title = "New Trip"
		}

		view.Render(w, r, "admin-trip", v)
	}

	var msg string

	t := models.Trip{
		ID:           utils.ToInt(f.ID),
		Title:        f.Title,
		Slug:         f.Slug,
		Status:       f.Status,
		Blurb:        f.Blurb,
		Description:  f.Description,
		Start:        domain.ToTime(f.Start),
		End:          domain.ToTime(f.End),
		TicketingURL: f.TicketingURL,
		Price:        f.Price,
		Notes:        f.Notes,
		ImageID:      f.ImageID,
		GalleryID:    f.GalleryID,
	}

	image, err := utils.UploadFile(w, r, "trip_image", "uploads/trip", true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(image) > 0 {
		t.ImageID = image[0].ID
	} else if (f.ImageID > 0) && (len(r.Form["deleteimg"]) == 1) {
		image := &models.File{
			ID: f.ImageID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		t.ImageID = 0
	}

	if t.ID != 0 {
		err := t.Update()
		if err != nil {
			if err == domain.ErrNotFound {
				view.NotFound(w, r)
				return
			}
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
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	id := strconv.Itoa(t.ID)

	http.Redirect(w, r, "/admin/trip?id="+id, http.StatusSeeOther)
}

func ListTrips(w http.ResponseWriter, r *http.Request) {
	trips, err := models.FetchTrips()
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

	err := t.GetBase()
	if err != nil {
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	if t.ImageID != 0 {
		image := &models.File{
			ID: t.ImageID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	}

	err = t.Delete()
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

	err := t.Fetch()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.FetchVendors(true)
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

	err := t.Fetch()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	vendors, err := models.FetchVendors(true)
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
