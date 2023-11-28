package app

import (
	"errors"
	"github.com/google/uuid"
	"html/template"
	"net/mail"
	"strings"
	"time"
)

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

func (m *Mail) SetRecipient(recipient string, domain string) error {
	m.To = recipient

	mailAddress, err := mail.ParseAddress(recipient)
	if err != nil {
		return err
	}

	if strings.HasSuffix(mailAddress.Address, domain) {
		return nil
	} else {
		return errors.New("mail recipient is to the wrong domain")
	}
}
