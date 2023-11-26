package smtp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"placemail/internal/util"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

type SmtpServer struct {
	Mail   map[string][]Mail
	Server *smtp.Server
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
				var newMail []Mail
				for _, mail := range mails {
					if time.Now().Sub(mail.Creation) <= time.Minute*10 {
						mail.Timestamp = util.GenerateTimestamp(mail.Creation)
						newMail = append(newMail, mail)
					}
				}

				s.Mail[email] = newMail
			}

			s.mu.Unlock()
		}
	}()
}

type Mail struct {
	From      string
	To        string
	Data      string
	Creation  time.Time
	Timestamp string
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
	mailServer.Mail = make(map[string][]Mail)

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
	}

	s.mail.Creation = time.Now()
	s.mail.Timestamp = util.GenerateTimestamp(s.mail.Creation)
	s.AppendMail()

	return nil
}

func (s *smtpSession) Reset() {}

func (s *smtpSession) Logout() error {
	return nil
}

func (s *smtpSession) AppendMail() {
	s.server.mu.Lock()
	defer s.server.mu.Unlock()

	if _, ok := s.server.Mail[s.mail.To]; !ok {
		s.server.Mail[s.mail.To] = make([]Mail, 0)
	}

	s.server.Mail[s.mail.To] = append(s.server.Mail[s.mail.To], s.mail)
}
