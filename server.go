package main

import (
	"github.com/joho/godotenv"
	"github.com/rutan/yarujan/lib/routes"
	"github.com/zenazn/goji"
	"gopkg.in/gographics/imagick.v1/imagick"
	"net/http"
)

func main() {
	godotenv.Load()

	imagick.Initialize()
	defer imagick.Terminate()

	goji.Get("/", http.FileServer(http.Dir("./public")))
	goji.Get("/assets/*", http.FileServer(http.Dir("./public")))
	goji.Get("/lgtm", routes.Index)
	goji.Post("/lgtm", routes.Create)
	goji.Serve()
}
