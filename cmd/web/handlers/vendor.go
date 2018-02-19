package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/flash"
	"strconv"
)

func VendorForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "vendor", &view.View{
			Form:  new(models.VendorForm),
			Title: "New Vendor",
		})
		return
	}

	v := &models.Vendor{
		ID: utils.ToInt(id),
	}

	err := v.Get()
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.VendorForm{
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

	view.Render(w, r, "vendor", &view.View{
		Title: f.Name,
		Form:  f,
	})
}

func PostVendor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.VendorForm{
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
		Active:  (len(r.Form["active"]) == 1),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if f.ID == "" {
			v.Title = "New Vendor"
		}

		view.Render(w, r, "vendor", v)
	}

	fn, err := utils.UploadFile(w, r, "brand_image", "uploads/vendor/")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(fn) > 0 && fn[0] != "" {
		f.Brand = fn[0]
	} else if (len(f.Brand) != 0) && (len(r.Form["deleteimg"]) == 1) {
		err = utils.DeleteFile(f.Brand)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		f.Brand = ""
	}

	var msg string

	v := models.Vendor{
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

	if v.ID != 0 {
		err := v.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		msg = utils.MsgSuccessfullyUpdated
	} else {
		err := v.Create()
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

	id := strconv.Itoa(v.ID)

	http.Redirect(w, r, "/admin/vendor?id="+id, http.StatusSeeOther)
}

func ListVendors(w http.ResponseWriter, r *http.Request) {
	vendors, err := models.GetVendors(false)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "vendors", &view.View{
		Title:   "Vendors",
		Vendors: vendors,
	})
}

func RemoveVendor(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	v := models.Vendor{
		ID: utils.ToInt(id),
	}

	err := v.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/vendors", http.StatusSeeOther)
}
