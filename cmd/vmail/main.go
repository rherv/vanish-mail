package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"vmail/internal/app"
	"vmail/internal/util"
)

//go:generate npm run build

var domain = flag.String("domain", "localhost", "the domain to accept emails for")
var httpPort = flag.Int("http", 8443, "http service address")
var smtpPort = flag.Int("smtp", 587, "smtp service address")
var delay = flag.Int("delay", 10, "the time in minutes to keep an email for")
var certFile = flag.String("cert", "server.crt", "Path to TLS certificate file")
var keyFile = flag.String("key", "server.key", "Path to TLS private key file")

func main() {
	flag.Parse()

	util.GenerateEmail(*domain)
	app.Init(*domain, *httpPort, *smtpPort, *delay, *certFile, *keyFile)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
}
