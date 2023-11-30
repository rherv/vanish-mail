package app

import (
	"fmt"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
)

func TestTrafficEmailServer(t *testing.T) {
	domain := "localhost"
	httpPort := 8080
	smtpPort := 1025
	delay := 10
	Init(domain, httpPort, smtpPort, delay)

	for i := 0; i < 100; i++ {
		SendEmail(domain, smtpPort, "test@mail.com", "testing@"+domain)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func SendEmail(domain string, port int, from string, to string) {
	subject := "Test Subject"
	body := "This is the body of the email."

	message := composeMessage(from, to, subject, body)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", domain, port), nil, from, []string{to}, []byte(message))

	if err != nil {
		return
	}
}

func composeMessage(from, to, subject, body string) string {
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-version": "1.0",
		"Content-Type": "text/plain; charset=\"UTF-8\"",
	}

	var messageBuilder strings.Builder

	for key, value := range headers {
		messageBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	messageBuilder.WriteString("\r\n" + body)

	return messageBuilder.String()
}
