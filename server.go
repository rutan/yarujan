package main

import (
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"gopkg.in/gographics/imagick.v2/imagick"
	"io/ioutil"
	"net/http"
)

func createLGTM(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	contents, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	mw.ReadImageBlob(contents)

	width := float64(mw.GetImageWidth())
	height := float64(mw.GetImageHeight())
	var size float64
	if width < height {
		size = width
	} else {
		size = height
	}
	max_size := float64(320)
	var ratio float64
	if size < max_size {
		ratio = 1
	} else {
		ratio = max_size / size
	}

	fmt.Printf("size: %d", size)
	err = mw.ResizeImage(uint(width*ratio), uint(height*ratio), imagick.FILTER_BOX, 1)
	if err != nil {
		http.Error(w, "crash image", http.StatusBadRequest)
		return
	}

	pw_fill := imagick.NewPixelWand()
	defer pw_fill.Destroy()
	pw_stroke := imagick.NewPixelWand()
	defer pw_stroke.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	dw.SetFont("Meiryo")
	dw.SetFontSize(width * ratio / 4)
	dw.SetTextAlignment(imagick.ALIGN_CENTER)
	dw.SetFontWeight(500)
	dw.SetStrokeWidth(3)
	pw_fill.SetColor("#ffffff")
	dw.SetFillColor(pw_fill)
	pw_stroke.SetColor("#444444")
	dw.SetStrokeColor(pw_stroke)
	dw.SetStrokeAntialias(true)
	mw.AnnotateImage(dw, width*ratio/2, height*ratio*0.95, 0, "LGTM")

	blob := mw.GetImageBlob()
	w.Write(blob)
}

func main() {
	goji.Post("/lgtm", createLGTM)
	goji.Get("/", http.FileServer(http.Dir("./public")))
	goji.Serve()
}
