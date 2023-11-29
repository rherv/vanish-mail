package app

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type App struct {
	SmtpServer    *EmailServer
	router        *mux.Router
	inboxTemplate *template.Template
	homeTemplate  *template.Template
	mailTemplate  *template.Template
	css           []byte
	Domain        string
	delay         int
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

//go:embed templates/inbox.html templates/home.html templates/mail.html templates/tailwind.css
var fileSystem embed.FS

func (a *App) loadTemplate(pattern string) *template.Template {
	tmpl, err := template.ParseFS(fileSystem, pattern)
	if err != nil {
		log.Fatalln(err)
	}

	return tmpl
}

func (a *App) templates() {
	a.inboxTemplate = a.loadTemplate("templates/inbox.html")
	a.homeTemplate = a.loadTemplate("templates/home.html")
	a.mailTemplate = a.loadTemplate("templates/mail.html")

	css, err := fileSystem.ReadFile("templates/tailwind.css")
	if err != nil {
		log.Fatalln(err)
	}

	a.css = css
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
			log.Fatalln(err)
			return
		}
	}()

	return &a
}
