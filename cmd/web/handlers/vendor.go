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
		BrandID: v.BrandID,
		Active:  v.Active,
	}

	if v.Brand != nil {
		f.Brand = v.Brand.Thumb
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
		BrandID: utils.ToInt(r.PostForm.Get("brand_id")),
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

	var msg string

	v := models.Vendor{
		ID:      utils.ToInt(f.ID),
		Name:    f.Name,
		Address: f.Address,
		City:    f.City,
		State:   f.State,
		Zip:     f.Zip,
		Phone:   f.Phone,
		Email:   f.Email,
		URL:     f.URL,
		Notes:   f.Notes,
		BrandID: f.BrandID,
		Active:  f.Active,
	}

	image, err := utils.UploadFile(w, r, "brand_image", "uploads/vendor", true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(image) > 0 {
		v.BrandID = image[0].ID
	} else if (f.BrandID > 0) && (len(r.Form["deleteimg"]) == 1) {
		image := &models.File{
			ID: f.BrandID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		v.BrandID = 0
	}

	if v.ID != 0 {
		err := v.Update()
		if err != nil {
			if err == db.ErrNotFound {
				view.NotFound(w, r)
				return
			}
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
	vars := mux.Vars(r)
	id := vars["id"]

	v := models.Vendor{
		ID: utils.ToInt(id),
	}

	err := v.GetBase()
	if err != nil {
		if err == db.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	if v.BrandID != 0 {
		image := &models.File{
			ID: v.BrandID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	}

	err = v.Delete()
	if err == db.ErrCannotDelete {
		err = flash.Add(w, r, utils.MsgCannotRemove, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	} else {
		err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "/admin/vendors", http.StatusSeeOther)
}
