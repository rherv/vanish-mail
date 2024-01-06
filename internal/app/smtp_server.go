package app

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
	"vmail/internal/util"

	"github.com/emersion/go-smtp"
)

type EmailServer struct {
	Mail       map[string]map[uuid.UUID]Mail
	SmtpServer *smtp.Server
	Delay      time.Duration
	mu         sync.RWMutex
}

func (s *EmailServer) Start() {
	s.RemoveOldMail()

	go func() {
		log.Println("Starting server at", s.SmtpServer.Addr)
		if err := s.SmtpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *EmailServer) RemoveOldMail() {
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

func NewSmtpServer(domain string, port int, delay int) *EmailServer {
	mailServer := &EmailServer{
		Delay: time.Duration(delay) * time.Minute,
	}

	s := smtp.NewServer(mailServer)

	s.Addr = fmt.Sprintf("%s:%d", domain, port)
	s.Domain = domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 50 * 1024 * 1024
	s.MaxLineLength = 2000
	s.MaxRecipients = 50
	s.AuthDisabled = false
	s.AllowInsecureAuth = true

	//s.EnableBINARYMIME = true
	// s.AuthDisabled = true
	// s.EnableSMTPUTF8 = true

	mailServer.SmtpServer = s
	mailServer.Mail = make(map[string]map[uuid.UUID]Mail)

	return mailServer
}
