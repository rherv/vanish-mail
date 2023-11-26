package main

import (
	"placemail/internal/util"
	"placemail/internal/web"
)

func main() {
	util.GenerateEmail("localhost")
	web.Init("localhost", 8080, 1025)
}
