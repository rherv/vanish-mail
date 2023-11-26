package main

import (
	"placemail/internal/web"
)

func main() {
	web.Init("localhost", 8080, 1025)
}
