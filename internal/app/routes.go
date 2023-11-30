package app

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"vmail/internal/util"
)

func (a *App) routes() {
	a.router.HandleFunc("/inbox/{email}/", a.inbox)
	a.router.HandleFunc("/inbox/{email}/{id}/", a.mail)
	a.router.HandleFunc("/inbox/{email}/{id}/delete/", a.delete)
	a.router.HandleFunc("/inbox/{email}/{id}/back/", a.back)
	a.router.HandleFunc("/tailwind.css", a.tailwind)
	a.router.HandleFunc("/", a.home)
}

func (a *App) tailwind(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	_, err := w.Write(a.css)
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *App) inbox(w http.ResponseWriter, r *http.Request) {
	var data pageData
	vars := mux.Vars(r)
	email := vars["email"]
	data.Email = email
	data.Delay = a.delay

	mail, ok := a.SmtpServer.Mail[email]
	var mails []Mail
	for _, value := range mail {
		mails = append(mails, value)
	}

	if ok {
		data.Mail = mails
	}

	err := a.inboxTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *App) mail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	mail, ok := a.SmtpServer.Mail[email][uid]
	if !ok {
		//http.NotFound(w, r)
		//return
	}

	data := pageData{
		Subject:   mail.Subject,
		Sender:    mail.From,
		Recipient: mail.To,
		Html:      mail.HTML,
	}

	err = a.mailTemplate.Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func (a *App) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	id := vars["id"]

	uid, err := uuid.Parse(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	delete(a.SmtpServer.Mail[email], uid)

	http.Redirect(w, r, fmt.Sprintf("/inbox/%s/", email), 302)
}

func (a *App) back(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	http.Redirect(w, r, fmt.Sprintf("/inbox/%s/", email), 302)
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	_ = r
	data := pageData{
		Email: util.GenerateEmail(a.Domain),
	}

	err := a.homeTemplate.Execute(w, data)

	if err != nil {
		log.Println(err)
		return
	}
}
