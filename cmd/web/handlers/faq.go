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

func FaqForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "faq-admin", &view.View{
			Form:  new(models.FAQForm),
			Title: "New FAQ",
		})
		return
	}

	faq := &models.FAQ{
		ID: utils.ToInt(id),
	}

	err := faq.Get()
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.FAQForm{
		ID:       strconv.Itoa(faq.ID),
		Question: faq.Question,
		Answer:   faq.Answer,
		Category: faq.Category,
		Order:    strconv.Itoa(faq.Order),
		Active:   faq.Active,
	}

	view.Render(w, r, "faq-admin", &view.View{
		Title: faq.Question,
		Form:  f,
		FAQ:   faq,
	})
}

func PostFAQ(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.FAQForm{
		ID:       r.PostForm.Get("id"),
		Question: r.PostForm.Get("question"),
		Category: r.PostForm.Get("category"),
		Answer:   r.PostForm.Get("answer"),
		Order:    r.PostForm.Get("order"),
		Active:   (id == "" || ((len(r.Form["active"]) == 1) && id != "")),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if id == "" {
			v.Title = "New FAQ"
		}

		view.Render(w, r, "faq-admin", v)
	}

	var msg string

	faq := models.FAQ{
		ID:       utils.ToInt(f.ID),
		Question: f.Question,
		Answer:   f.Answer,
		Category: f.Category,
		Order:    utils.ToInt(f.Order),
		Active:   f.Active,
	}

	if id != "" {
		faq.ID = utils.ToInt(id)

		err := faq.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		msg = utils.MsgSuccessfullyUpdated
	} else {
		err := faq.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		id = strconv.Itoa(faq.ID)
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/faq?id="+id, http.StatusSeeOther)
}

func ListFAQs(w http.ResponseWriter, r *http.Request) {
	faqs, err := models.GetFAQs()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "faqs-admin", &view.View{
		Title: "FAQs",
		FAQs:  faqs,
	})
}

func RemoveFAQ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	faq := models.FAQ{
		ID: utils.ToInt(id),
	}

	err := faq.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/faqs", http.StatusSeeOther)
}
