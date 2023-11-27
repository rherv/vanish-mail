package app

import (
	"embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

type App struct {
	SmtpServer    *SmtpServer
	router        *mux.Router
	inboxTemplate *template.Template
	homeTemplate  *template.Template
	mailTemplate  *template.Template
	Domain        string
	delay         int
}

type Mail struct {
	UUID      uuid.UUID
	Subject   string
	From      string
	To        string
	Timestamp string
	Data      []byte
	Creation  time.Time
	HTML      template.HTML
}

type pageData struct {
	Email     string
	Mail      []Mail
	Delay     int
	Subject   string
	Sender    string
	Recipient string
	Html      template.HTML
}

//go:embed templates/inbox.html
var inboxTemplate embed.FS

//go:embed templates/home.html
var homeTemplate embed.FS

//go:embed templates/mail.html
var mailTemplate embed.FS

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

	tmpl, err = template.ParseFS(mailTemplate, "templates/mail.html")
	if err != nil {
		log.Fatalln(err)
	}
	a.mailTemplate = tmpl
}

func Init(domain string, httpPort int, mailPort int, delay int) *App {
	a := App{
		SmtpServer: NewSmtpServer(domain, mailPort, delay),
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
