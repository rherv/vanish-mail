package app

import (
	"errors"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"io"
	"placemail/internal/util"
	"time"
)

type SmtpSession struct {
	mail   Mail
	server *SmtpServer
}

func (s *SmtpServer) NewSession(conn *smtp.Conn) (smtp.Session, error) {
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
	s.mail.From = from
	return nil
}

func (s *SmtpSession) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.mail.To = to

	if !util.IsValid(s.mail.To, s.server.Server.Domain) {
		return errors.New("the received email is to the wrong domain: ")
	}

	return nil
}

func (s *SmtpSession) Data(r io.Reader) error {
	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		subject, html, err := util.ParseEmailData(data)
		if err != nil {
			return err
		}

		s.mail.Subject = subject
		s.mail.HTML = html
		s.mail.UUID = uuid.New()
	}

	s.mail.Creation = time.Now()
	s.mail.Timestamp = util.GenerateTimestamp(s.mail.Creation)
	s.AppendMail()

	return nil
}

func (s *SmtpSession) AppendMail() {
	s.server.mu.Lock()
	defer s.server.mu.Unlock()

	if _, ok := s.server.Mail[s.mail.To]; !ok {
		s.server.Mail[s.mail.To] = make(map[uuid.UUID]Mail)
	}

	s.server.Mail[s.mail.To][s.mail.UUID] = s.mail
}

func (s *SmtpSession) Reset() {}

func (s *SmtpSession) Logout() error {
	return nil
}
