package app

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
)

func TestTrafficEmailServer(t *testing.T) {
	domain := "localhost"
	httpPort := 8080
	smtpPort := 1025
	delay := 10
	_ = Init(domain, httpPort, smtpPort, delay)

	for i := 0; i < 10; i++ {
		SendEmail(domain, smtpPort, strconv.Itoa(i)+"test@mail.com", "testing@"+domain)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
}

func SendEmail(domain string, port int, from string, to string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", domain, port))
	if err != nil {
		return
	}

	var data []string
	data = append(data, fmt.Sprintf("EHLO %s\r\n", domain))
	data = append(data, fmt.Sprintf("MAIL FROM:<%s>\r\n", from))
	data = append(data, fmt.Sprintf("RCPT TO:<%s>\r\n", to))
	data = append(data, "DATA\r\n")
	data = append(data, "Subject: Hello World\r\n")
	data = append(data, fmt.Sprintf("From: <%s>\r\n", from))
	data = append(data, fmt.Sprintf("To: <%s>\r\n", to))
	data = append(data, "Content-Type: text/html\r\n")
	data = append(data, `
<html>
	<body>
		<h1> hello! </h1>
		<p>This is your HTML content.</p>
	</body>
</html>`+"\r\n")

	data = append(data, ".\r\n")

	for _, line := range data {
		_, err := fmt.Fprintf(conn, line)
		if err != nil {
			return
		}
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
