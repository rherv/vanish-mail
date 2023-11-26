package web

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"placemail/internal/smtp"
)

type app struct {
	smtpServer    *smtp.SmtpServer
	router        *mux.Router
	inboxTemplate *template.Template
}

type inboxData struct {
	Email string
	Mail  []smtp.Mail
}

//go:embed templates/inbox.html
var inboxTemplate embed.FS

func (a *app) inbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data inboxData
	email := vars["email"]

	mail, ok := a.smtpServer.Mail[email]
	if !ok {
		data = inboxData{
			Email: email,
		}
	} else {
		data = inboxData{
			Email: email,
			Mail:  mail,
		}
	}

	err := a.inboxTemplate.Execute(w, data)
	if err != nil {
		return
	}
}

func (a *app) templates() {
	tmpl, err := template.ParseFS(inboxTemplate, "templates/inbox.html")
	if err != nil {
		log.Println("error here wtf")
		log.Fatalln(err)
	}

	a.inboxTemplate = tmpl
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
	a.templates()

	addr := fmt.Sprintf("%s:%d", domain, httpPort)

	err := http.ListenAndServe(addr, a.router)
	if err != nil {
		return
	}
}
