package handlers

import (
	"database/sql"
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/domain/models"
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

	err := v.Fetch()
	if err != nil {
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	f := &models.VendorForm{
		ID:      strconv.Itoa(v.ID),
		Name:    v.Name.String,
		Address: v.Address.String,
		City:    v.City.String,
		State:   v.State.String,
		Zip:     v.Zip.String,
		Phone:   v.Phone.String,
		Email:   v.Email.String,
		URL:     v.URL.String,
		Notes:   v.Notes.String,
		BrandID: int(v.BrandID.Int64),
		Active:  v.Active,
	}

	if v.Brand != nil {
		f.Brand = v.Brand.Thumb.String
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
		Name:    utils.NewNullStr(f.Name),
		Address: utils.NewNullStr(f.Address),
		City:    utils.NewNullStr(f.City),
		State:   utils.NewNullStr(f.State),
		Zip:     utils.NewNullStr(f.Zip),
		Phone:   utils.NewNullStr(f.Phone),
		Email:   utils.NewNullStr(f.Email),
		URL:     utils.NewNullStr(f.URL),
		Notes:   utils.NewNullStr(f.Notes),
		Active:  f.Active,
	}

	if f.BrandID != 0 {
		v.BrandID = utils.NewNullInt(f.BrandID)
	} else {
		v.BrandID = sql.NullInt64{}
	}

	image, err := utils.UploadFile(w, r, "brand_image", "uploads/vendor", true)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if len(image) > 0 {
		v.BrandID = utils.NewNullInt(image[0].ID)
	} else if (f.BrandID != 0) && (len(r.Form["deleteimg"]) == 1) {
		image := &models.File{
			ID: f.BrandID,
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		v.BrandID = sql.NullInt64{}
	}

	if v.ID != 0 {
		err := v.Update()
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
	vendors, err := models.FetchVendors(false)
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
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	if int(v.BrandID.Int64) != 0 {
		image := &models.File{
			ID: int(v.BrandID.Int64),
		}

		err = utils.DeleteFile(image)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	}

	err = v.Delete()
	if err != nil {
		if err == domain.ErrCannotDelete {
			err = flash.Add(w, r, utils.MsgCannotRemove, "warning")
			if err != nil {
				view.ServerError(w, r, err)
				return
			}
		}
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
