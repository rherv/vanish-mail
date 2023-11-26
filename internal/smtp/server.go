package smtp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

type SmtpServer struct {
	MailChannel chan Mail
	Server      *smtp.Server
}

func (s *SmtpServer) Start() {
	go func() {
		log.Println("Starting server at", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

type Mail struct {
	From string
	To   string
	Data string
}

func NewSmtpServer(domain string, port int) *SmtpServer {
	mailServer := &SmtpServer{}

	s := smtp.NewServer(mailServer)

	s.Addr = fmt.Sprintf("%s:%d", domain, port)
	s.Domain = domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	mailServer.Server = s
	mailServer.MailChannel = make(chan Mail)

	return mailServer
}

type smtpSession struct {
	mail   Mail
	server *SmtpServer
}

func (s *SmtpServer) NewSession(c *smtp.Conn) (smtp.Session, error) {
	_ = c

	return &smtpSession{
		mail:   Mail{},
		server: s,
	}, nil
}

func (s *smtpSession) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}
	return nil
}

func (s *smtpSession) Mail(from string, opts *smtp.MailOptions) error {
	_ = opts

	s.mail.From = from
	return nil
}

func (s *smtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	_ = opts

	s.mail.To = to
	return nil
}

func (s *smtpSession) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		s.mail.Data = string(b)
		s.server.MailChannel <- s.mail
	}
	return nil
}

func (s *smtpSession) Reset() {}

func (s *smtpSession) Logout() error {
	return nil
}
