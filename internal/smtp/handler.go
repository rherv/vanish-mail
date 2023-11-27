package smtp

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"log"
	"placemail/internal/util"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

type SmtpServer struct {
	Mail   map[string]map[uuid.UUID]Mail
	Server *smtp.Server
	Delay  time.Duration
	mu     sync.RWMutex
}

func (s *SmtpServer) Start() {
	s.RemoveOldMail()

	go func() {
		log.Println("Starting server at", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *SmtpServer) RemoveOldMail() {
	go func() {
		for {
			time.Sleep(time.Second * 10)

			s.mu.Lock()

			for email, mails := range s.Mail {
				for id, mail := range mails {
					m := s.Mail[email][id]
					m.Timestamp = util.GenerateTimestamp(mail.Creation)
					if time.Now().Sub(mail.Creation) >= s.Delay {
						delete(mails, id)
					} else {
						s.Mail[email][id] = m
					}
				}
			}

			s.mu.Unlock()
		}
	}()
}

type Mail struct {
	UUID      uuid.UUID
	Subject   string
	From      string
	To        string
	Data      string
	Creation  time.Time
	Timestamp string
	HTML      template.HTML
}

func NewSmtpServer(domain string, port int, delay int) *SmtpServer {
	mailServer := &SmtpServer{
		Delay: time.Duration(delay) * time.Minute,
	}

	s := smtp.NewServer(mailServer)

	s.Addr = fmt.Sprintf("%s:%d", domain, port)
	s.Domain = domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true
	s.EnableSMTPUTF8 = true
	s.EnableBINARYMIME = true
	mailServer.Server = s
	mailServer.Mail = make(map[string]map[uuid.UUID]Mail)

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

	if !util.IsValid(s.mail.To, s.server.Server.Domain) {
		return errors.New("the received email is to the wrong domain: ")
	}

	return nil
}

func (s *smtpSession) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		log.Println(err)
		return err
	} else {
		html, err := util.ParseEmailData(b)
		if err != nil {
			return err
		}
		s.mail.HTML = html

		s.mail.Data = string(b)

		s.mail.UUID = uuid.New()
	}

	s.mail.Creation = time.Now()
	s.mail.Timestamp = util.GenerateTimestamp(s.mail.Creation)
	s.AppendMail()

	return nil
}

func (s *smtpSession) AppendMail() {
	s.server.mu.Lock()
	defer s.server.mu.Unlock()

	if _, ok := s.server.Mail[s.mail.To]; !ok {
		s.server.Mail[s.mail.To] = make(map[uuid.UUID]Mail)

		s.server.Mail[s.mail.To][s.mail.UUID] = s.mail
	}
}

func (s *smtpSession) Reset() {}

func (s *smtpSession) Logout() error {
	return nil
}
