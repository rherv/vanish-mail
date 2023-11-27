package util

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/jhillyerd/enmime"
	"html/template"
	"log"
	"math/rand"
	"net/mail"
	"strings"
)

func ParseEmailData(data []byte) (subject string, html template.HTML, err error) {
	envelope, err := enmime.ReadEnvelope(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return envelope.GetHeader("Subject"), template.HTML(envelope.HTML), nil
}

//go:embed first-names.txt
var nameList string

var names = strings.Split(strings.TrimSpace(nameList), "\n")

func GenerateEmail(domain string) string {
	first := names[rand.Intn(len(names))]
	last := names[rand.Intn(len(names))]

	first = strings.TrimSuffix(first, "\n")
	last = strings.TrimSuffix(last, "\n")

	email := fmt.Sprintf("%s.%s@%s", first, last, domain)

	return email
}

func IsValid(email string, domain string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return strings.HasSuffix(addr.Address, domain)
}
