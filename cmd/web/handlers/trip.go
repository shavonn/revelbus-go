package handlers

import (
	"database/sql"
	"net/http"
	"revelbus/cmd/web/utils"
	"revelbus/cmd/web/view"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/domain/models"
	"revelbus/internal/platform/flash"
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
		Title:        t.Title.String,
		Slug:         t.Slug.String,
		Status:       t.Status.String,
		Blurb:        t.Blurb.String,
		Description:  t.Description.String,
		Start:        t.Start.Format(domain.TimeFormat),
		End:          t.End.Format(domain.TimeFormat),
		Price:        t.Price.String,
		TicketingURL: t.TicketingURL.String,
		Notes:        t.Notes.String,
		ImageID:      int(t.ImageID.Int64),
		GalleryID:    int(t.GalleryID.Int64),
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
		Title:        utils.NewNullStr(f.Title),
		Slug:         utils.NewNullStr(f.Slug),
		Status:       utils.NewNullStr(f.Status),
		Blurb:        utils.NewNullStr(f.Blurb),
		Description:  utils.NewNullStr(f.Description),
		Start:        domain.ToTime(f.Start),
		End:          domain.ToTime(f.End),
		TicketingURL: utils.NewNullStr(f.TicketingURL),
		Price:        utils.NewNullStr(f.Price),
		Notes:        utils.NewNullStr(f.Notes),
	}

	if f.ImageID != 0 {
		t.ImageID = utils.NewNullInt(f.ImageID)
	} else {
		t.ImageID = sql.NullInt64{}
	}

	if f.GalleryID != 0 {
		t.GalleryID = utils.NewNullInt(f.GalleryID)
	} else {
		t.GalleryID = sql.NullInt64{}
	}

	image, err := utils.UploadFile(w, r, "trip_image", "uploads/trip", true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(image) > 0 {
		t.ImageID = utils.NewNullInt(image[0].ID)
	} else if (f.ImageID != 0) && (len(r.Form["deleteimg"]) == 1) {
		image := &models.File{
			ID: f.ImageID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		t.ImageID = sql.NullInt64{}
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

	if int(t.ImageID.Int64) != 0 {
		image := &models.File{
			ID: int(t.ImageID.Int64),
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
