package main

import (
	"./lib/routes"
	"github.com/joho/godotenv"
	"github.com/zenazn/goji"
	"gopkg.in/gographics/imagick.v2/imagick"
	"net/http"
)

func main() {
	godotenv.Load()

	imagick.Initialize()
	defer imagick.Terminate()

	goji.Get("/", http.FileServer(http.Dir("./public")))
	goji.Get("/assets/*", http.FileServer(http.Dir("./public")))
	goji.Post("/lgtm", routes.Create)
	goji.Serve()
}
