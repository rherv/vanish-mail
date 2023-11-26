package web

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"placemail/internal/smtp"
	"placemail/internal/util"
)

type app struct {
	smtpServer    *smtp.SmtpServer
	router        *mux.Router
	inboxTemplate *template.Template
	homeTemplate  *template.Template
	domain        string
}

type pageData struct {
	Email string
	Mail  []smtp.Mail
}

//go:embed templates/inbox.html
var inboxTemplate embed.FS

//go:embed templates/home.html
var homeTemplate embed.FS

func (a *app) inbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data pageData
	email := vars["email"]
	data.Email = email

	mail, ok := a.smtpServer.Mail[email]
	if ok {
		data.Mail = mail
	}

	err := a.inboxTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Email: util.GenerateEmail(a.domain),
	}

	err := a.homeTemplate.Execute(w, data)

	if err != nil {
		log.Println(err)
		return
	}
}

func (a *app) templates() {
	tmpl, err := template.ParseFS(inboxTemplate, "templates/inbox.html")
	if err != nil {
		log.Fatalln(err)
	}
	a.inboxTemplate = tmpl

	tmpl, err = template.ParseFS(homeTemplate, "templates/home.html")
	if err != nil {
		log.Fatalln(err)
	}
	a.homeTemplate = tmpl
}

func (a *app) routes() {
	a.router.HandleFunc("/inbox/{email}", a.inbox)
	a.router.HandleFunc("/", a.home)
}

func Init(domain string, httpPort int, mailPort int) {
	a := app{
		smtpServer: smtp.NewSmtpServer(domain, mailPort),
		router:     mux.NewRouter(),
	}

	a.domain = domain
	a.smtpServer.Start()
	a.routes()
	a.templates()

	addr := fmt.Sprintf("%s:%d", domain, httpPort)

	err := http.ListenAndServe(addr, a.router)
	if err != nil {
		return
	}
}
