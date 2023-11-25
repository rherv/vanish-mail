package mailserver

import (
	"fmt"
	"net"
	"testing"
)

func TestSmtpServer(t *testing.T) {
	domain := "localhost"
	port := 1025

	sender := "sender@test.com"
	recipient := "recipient@test.com"
	data := "Hello World!"

	mailServer := NewSmtpServer(domain, port)
	go mailServer.Start()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", domain, port))
	if err != nil {
		t.Error(err)
	}

	err = nil
	_, err = fmt.Fprintf(conn, "EHLO localhost\r\n")
	_, err = fmt.Fprintf(conn, "AUTH PLAIN\r\n")
	_, err = fmt.Fprintf(conn, "AHVzZXJuYW1lAHBhc3N3b3Jk\r\n")
	_, err = fmt.Fprintf(conn, "MAIL FROM:%s\r\n", sender)
	_, err = fmt.Fprintf(conn, "RCPT TO:%s\r\n", recipient)
	_, err = fmt.Fprintf(conn, "DATA\r\n")
	_, err = fmt.Fprintf(conn, "%s\r\n", data)
	_, err = fmt.Fprintf(conn, ".\r\n")
	if err != nil {
		t.Error(err)
	}

	mail, ok := <-mailServer.MailChannel
	if !ok {
		t.Error("failed to read mail channel")
	}

	if mail.From != sender || mail.To != recipient || mail.Data != data+"\r\n" {
		t.Errorf("invalid mail data")
		t.Errorf("expected: %s %s %s", sender, recipient, data)
		t.Errorf("actual: %s %s %s", mail.From, mail.To, mail.Data)
	}
}