package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
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

	v := &db.Vendor{
		ID: toInt(id),
	}

	err := v.Get()
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
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
		Brand:   v.Brand,
		Active:  v.Active,
	}

	render(w, r, "vendor.html", &view{
		Title:  f.Name,
		Form:   f,
		Vendor: v,
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
		Brand:   r.PostForm.Get("brand"),
		Active:  (id == "" || ((len(r.Form["active"]) == 1) && id != "")),
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

	fn, err := uploadFile(w, r, "brand_image", "uploads/vendor/")
	if err != nil {
		serverError(w, r, err)
		return
	}

	if fn != "" {
		f.Brand = fn
	} else if (len(f.Brand) != 0) && (len(r.Form["deleteimg"]) == 1) {
		err = deleteFile("uploads/vendor/" + f.Brand)
		if err != nil {
			serverError(w, r, err)
			return
		}
		f.Brand = ""
	}

	var msg string

	v := db.Vendor{
		Name:    f.Name,
		Address: f.Address,
		City:    f.City,
		State:   f.State,
		Zip:     f.Zip,
		Phone:   f.Phone,
		Email:   f.Email,
		URL:     f.URL,
		Notes:   f.Notes,
		Brand:   f.Brand,
		Active:  f.Active,
	}

	if id != "" {
		v.ID = toInt(id)

		err := v.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		msg = MsgSuccessfullyUpdated
	} else {
		err := v.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		id = strconv.Itoa(v.ID)
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
	vendors, err := db.GetVendors(false)
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "vendors.html", &view{
		Title:   "Vendors",
		Vendors: &vendors,
	})
}

func removeVendor(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	v := db.Vendor{
		ID: toInt(id),
	}

	err := v.Delete()
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
