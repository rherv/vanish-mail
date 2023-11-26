package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"placemail/internal/smtp"
)

type app struct {
	smtpServer *smtp.SmtpServer
	router     *mux.Router
}

func (a *app) inbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]

	_, err := fmt.Fprintf(w, email)
	if err != nil {
		return
	}
}

func (a *app) routes() {
	a.router.HandleFunc("/inbox/{email}", a.inbox)
}

func Init(domain string, httpPort int, mailPort int) {
	a := app{
		smtpServer: smtp.NewSmtpServer(domain, mailPort),
		router:     mux.NewRouter(),
	}

	a.smtpServer.Start()
	a.routes()

	addr := fmt.Sprintf("%s:%d", domain, httpPort)

	err := http.ListenAndServe(addr, a.router)
	if err != nil {
		return
	}
}
