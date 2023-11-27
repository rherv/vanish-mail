package main

import (
	"flag"
	"placemail/internal/util"
	"placemail/internal/web"
)

var domain = flag.String("domain", "localhost", "the domain to accept emails for")
var httpPort = flag.Int("http", 8080, "http service address")
var smtpPort = flag.Int("smtp", 1025, "smtp service address")
var delay = flag.Int("delay", 10, "the time in minutes to keep an email for")

func main() {
	flag.Parse()

	util.GenerateEmail(*domain)
	web.Init(*domain, *httpPort, *smtpPort, *delay)
}
