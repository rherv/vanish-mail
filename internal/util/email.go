package util

import (
	_ "embed"
	"fmt"
	"math/rand"
	"net/mail"
	"strings"
)

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
