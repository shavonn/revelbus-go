package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strconv"

	"github.com/gorilla/mux"
)

func faqForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		render(w, r, "faq-admin", &view{
			Form:  new(forms.FAQForm),
			Title: "New FAQ",
		})
		return
	}

	faq := &db.FAQ{
		ID: toInt(id),
	}

	err := faq.Get()
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.FAQForm{
		ID:       strconv.Itoa(faq.ID),
		Question: faq.Question,
		Answer:   faq.Answer,
		Category: faq.Category,
		Order:    strconv.Itoa(faq.Order),
		Active:   faq.Active,
	}

	render(w, r, "faq-admin", &view{
		Title: faq.Question,
		Form:  f,
		FAQ:   faq,
	})
}

func postFAQ(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.FAQForm{
		ID:       r.PostForm.Get("id"),
		Question: r.PostForm.Get("question"),
		Category: r.PostForm.Get("category"),
		Answer:   r.PostForm.Get("answer"),
		Order:    r.PostForm.Get("order"),
		Active:   (id == "" || ((len(r.Form["active"]) == 1) && id != "")),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		if id == "" {
			v.Title = "New FAQ"
		}

		render(w, r, "faq-admin", v)
	}

	var msg string

	faq := db.FAQ{
		ID:       toInt(f.ID),
		Question: f.Question,
		Answer:   f.Answer,
		Category: f.Category,
		Order:    toInt(f.Order),
		Active:   f.Active,
	}

	if id != "" {
		faq.ID = toInt(id)

		err := faq.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		msg = MsgSuccessfullyUpdated
	} else {
		err := faq.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		id = strconv.Itoa(faq.ID)
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/faq?id="+id, http.StatusSeeOther)
}

func listFAQs(w http.ResponseWriter, r *http.Request) {
	faqs, err := db.GetFAQs()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "faqs-admin", &view{
		Title: "FAQs",
		FAQs:  faqs,
	})
}

func removeFAQ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	faq := db.FAQ{
		ID: toInt(id),
	}

	err := faq.Delete()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyRemoved, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/faqs", http.StatusSeeOther)
}
