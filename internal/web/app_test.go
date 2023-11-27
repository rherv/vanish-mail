package web

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTrafficSmtpServer(t *testing.T) {
	a := Init("localhost", 8080, 1025, 10)

	conn, err := net.Dial("tcp", "localhost:1025")
	if err != nil {
		return
	}

	var data []string
	data = append(data, "EHLO localhost\r\n")
	data = append(data, "MAIL FROM:<root@nsa.gov>\r\n")
	data = append(data, "RCPT TO:<testing@localhost>\r\n")
	data = append(data, "DATA\r\n")
	data = append(data, "Subject: Hello World\r\n")
	data = append(data, "From: <testing@localhost>\r\n")
	data = append(data, "To: <recipient_email@example.com>\r\n")
	data = append(data, "Content-Type: text/html\r\n")
	data = append(data, `
<html>
	<body>
		<p>This is your HTML content.</p>
	</body>
</html>`+"\r\n")

	data = append(data, ".\r\n")

	for _, line := range data {
		log.Print(line)
		_, err := fmt.Fprintf(conn, line)
		if err != nil {
			return
		}
	}

	_ = a

	b, err := io.ReadAll(conn)
	if err != nil {
		return
	}

	log.Println("response:", string(b))

	time.Sleep(1 * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
