package handlers

import (
	"net/http"
	"revelforce-admin/internal/platform/db"
	"revelforce-admin/internal/platform/flash"
	"revelforce-admin/internal/platform/forms"
	"strconv"
)

func vendorForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		render(w, r, "vendor.html", &view{
			Form:  new(forms.VendorForm),
			Title: "New Vendor",
		})
		return
	}

	v, err := db.GetVendorByID(id)
	if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.VendorForm{
		ID:      strconv.Itoa(v.ID),
		Name:    v.Name,
		Address: v.Address,
		City:    v.City,
		State:   v.State,
		Zip:     v.Zip,
		Phone:   v.Phone,
		Email:   v.Email,
		URL:     v.URL,
		Notes:   v.Notes,
	}

	render(w, r, "vendor.html", &view{
		Form:   f,
		Vendor: *v,
	})
}

func postVendor(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.VendorForm{
		ID:      r.PostForm.Get("id"),
		Name:    r.PostForm.Get("name"),
		Address: r.PostForm.Get("address"),
		City:    r.PostForm.Get("city"),
		State:   r.PostForm.Get("state"),
		Zip:     r.PostForm.Get("zip"),
		Phone:   r.PostForm.Get("phone"),
		Email:   r.PostForm.Get("email"),
		URL:     r.PostForm.Get("url"),
		Notes:   r.PostForm.Get("notes"),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		if id == "" {
			v.Title = "New Vendor"
		}

		render(w, r, "vendor.html", v)
	}

	var msg string

	vendor := db.Vendor{
		Name:    f.Name,
		Address: f.Address,
		City:    f.City,
		State:   f.State,
		Zip:     f.Zip,
		Phone:   f.Phone,
		Email:   f.Email,
		URL:     f.URL,
		Notes:   f.Notes,
	}

	if id != "" {
		intID, _ := strconv.Atoi(id)
		vendor.ID = intID

		err := vendor.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		msg = MsgSuccessfullyUpdated
	} else {
		vid, err := vendor.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		id = strconv.Itoa(vid)
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/vendor?id="+id, http.StatusSeeOther)
}

func listVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := db.GetVendors()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "vendors.html", &view{
		Title:   "Vendors",
		Vendors: vendors,
	})
}

func removeVendor(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	v, err := db.GetVendorByID(id)
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = v.Delete()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyRemoved, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/vendors", http.StatusSeeOther)
}
