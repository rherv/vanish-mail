package app

import (
	"errors"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"html/template"
	"io"
	"sync"
	"time"
	"vmail/internal/util"
)

type SmtpSession struct {
	mail   Mail
	server *EmailServer
	mu     sync.Mutex
}

func (s *EmailServer) NewSession(conn *smtp.Conn) (smtp.Session, error) {
	return &SmtpSession{
		mail:   Mail{},
		server: s,
	}, nil
}

func (s *SmtpSession) AuthPlain(username, password string) error {
	/*
		if username != "username" || password != "password" {
			return errors.New("invalid username or password")
		}
	*/

	return nil
}

func (s *SmtpSession) Mail(from string, opts *smtp.MailOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mail.From = from
	return nil
}

func (s *SmtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.mail.SetRecipient(to, s.server.SmtpServer.Domain)
	if err != nil {
		return err
	}

	return nil
}

func (s *SmtpSession) Data(r io.Reader) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	envelope, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	} else {
		contentType := envelope.GetHeader("Content-Type")

		if util.IsHtml(contentType) {
			s.mail.HTML = template.HTML(envelope.HTML)
		} else if util.IsText(contentType) {
			s.mail.HTML = template.HTML(envelope.Text)
		} else {
			return errors.New("invalid content-type")
		}

		s.mail.Subject = envelope.GetHeader("Subject")
	}

	s.mail.UUID = uuid.New()
	s.mail.Creation = time.Now()
	s.mail.Timestamp = util.GenerateTimestamp(s.mail.Creation)
	s.AppendMail()

	return nil
}

func (s *SmtpSession) AppendMail() {
	if _, ok := s.server.Mail[s.mail.To]; !ok {
		s.server.Mail[s.mail.To] = make(map[uuid.UUID]Mail)
	}

	s.server.Mail[s.mail.To][s.mail.UUID] = s.mail
}

func (s *SmtpSession) Reset() {}

func (s *SmtpSession) Logout() error {
	return nil
}
