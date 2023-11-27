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

type App struct {
	SmtpServer    *smtp.SmtpServer
	router        *mux.Router
	inboxTemplate *template.Template
	homeTemplate  *template.Template
	Domain        string
	delay         int
}

type pageData struct {
	Email string
	Mail  []smtp.Mail
	Delay int
}

//go:embed templates/inbox.html
var inboxTemplate embed.FS

//go:embed templates/home.html
var homeTemplate embed.FS

func (a *App) inbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data pageData
	email := vars["email"]
	data.Email = email
	data.Delay = a.delay

	mail, ok := a.SmtpServer.Mail[email]
	if ok {
		data.Mail = mail
	}

	err := a.inboxTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Email: util.GenerateEmail(a.Domain),
	}

	err := a.homeTemplate.Execute(w, data)

	if err != nil {
		log.Println(err)
		return
	}
}

func (a *App) templates() {
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

func (a *App) routes() {
	a.router.HandleFunc("/inbox/{email}", a.inbox)
	a.router.HandleFunc("/", a.home)
}

func Init(domain string, httpPort int, mailPort int, delay int) *App {
	a := App{
		SmtpServer: smtp.NewSmtpServer(domain, mailPort, delay),
		router:     mux.NewRouter(),
		delay:      delay,
	}

	a.Domain = domain
	a.SmtpServer.Start()
	a.routes()
	a.templates()

	addr := fmt.Sprintf("%s:%d", domain, httpPort)

	go func() {
		err := http.ListenAndServe(addr, a.router)
		if err != nil {
			return
		}
	}()

	return &a
}
