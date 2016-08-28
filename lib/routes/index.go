package routes

import (
	"../uploader"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"
	"os"
)

type ImageIndex struct {
	Images []Image `json:"images"`
}

func Index(c web.C, w http.ResponseWriter, r *http.Request) {
	client := uploader.New()
	list, err := client.GetURLList(os.Getenv("AWS_S3_BUCKET_NAME"))
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	images := make([]Image, len(list))
	for i, url := range list {
		images[i] = Image{
			Url: url,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(ImageIndex{Images: images})
}
