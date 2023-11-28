package app

import (
	"bytes"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"html/template"
	"io"
	"log"
	"placemail/internal/util"
	"sync"
	"time"
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

	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		//subject, html, err := util.ParseEmailData(data)
		if err != nil {
			return err
		}

		err := s.ParseData(data)
		if err != nil {
			return err
		}
	}

	s.mail.UUID = uuid.New()
	s.mail.Creation = time.Now()
	s.mail.Timestamp = util.GenerateTimestamp(s.mail.Creation)
	s.AppendMail()

	return nil
}

func (s *SmtpSession) ParseData(data []byte) error {
	envelope, err := enmime.ReadEnvelope(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return err
	}

	s.mail.Subject = envelope.GetHeader("Subject")
	s.mail.HTML = template.HTML(envelope.HTML)

	return nil
}

func (s *SmtpSession) AppendMail() {
	s.mu.Lock()
	defer s.mu.Unlock()

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
