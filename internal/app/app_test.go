package app

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTrafficSmtpServer(t *testing.T) {
	domain := "localhost"
	httpPort := 8080
	smtpPort := 1025
	delay := 10
	a := Init(domain, httpPort, smtpPort, delay)

	start := time.Now()

	for i := 0; i < 10; i++ {
		SendEmail(domain, smtpPort, strconv.Itoa(i)+"test@mail.com", "testing@"+domain)
	}

	log.Println("elapsed:", time.Now().Sub(start))

	_ = a

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
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
