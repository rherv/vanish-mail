package app

import (
	"crypto/tls"
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
		if err := s.SmtpServer.ListenAndServeTLS(); err != nil {
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

func NewSmtpServer(domain string, port int, delay int, certFile string, keyFile string) *EmailServer {
	mailServer := &EmailServer{
		Delay: time.Duration(delay) * time.Minute,
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error loading certificate:", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	s := smtp.NewServer(mailServer)

	s.Addr = fmt.Sprintf("%s:%d", domain, port)
	s.Domain = domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024 * 1024
	s.MaxLineLength = 2000
	s.MaxRecipients = 50
	s.EnableREQUIRETLS = true
	s.EnableDSN = true
	s.EnableSMTPUTF8 = true
	s.AllowInsecureAuth = false
	s.AuthDisabled = true
	s.ErrorLog = &log.Logger{}

	//s.EnableBINARYMIME = true
	// s.AuthDisabled = true
	// s.EnableSMTPUTF8 = true
	s.TLSConfig = tlsConfig

	mailServer.SmtpServer = s
	mailServer.Mail = make(map[string]map[uuid.UUID]Mail)

	return mailServer
}
