package main

import (
	"log"
	"placemail/internal/util"
	"placemail/internal/web"
)

func main() {
	log.Println("is valid", util.IsValid("test@localhost", "localhost"))
	util.GenerateEmail("localhost")
	web.Init("localhost", 8080, 1025)
}
