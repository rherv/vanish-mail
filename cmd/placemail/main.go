package main

import "placemail/internal/mailserver"

func main() {
	mailServer := mailserver.NewSmtpServer("localhost", 1025)
	mailServer.Start()
}
