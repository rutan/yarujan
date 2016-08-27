package main

import (
	"./lib/routes"
	"github.com/zenazn/goji"
	"gopkg.in/gographics/imagick.v2/imagick"
	"net/http"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	goji.Post("/lgtm", routes.Create)
	goji.Get("/", http.FileServer(http.Dir("./public")))
	goji.Serve()
}
