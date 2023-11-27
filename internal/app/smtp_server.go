package app

import (
	"fmt"
	"github.com/google/uuid"
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
	mailServer.Server = s
	mailServer.Mail = make(map[string]map[uuid.UUID]Mail)

	return mailServer
}
