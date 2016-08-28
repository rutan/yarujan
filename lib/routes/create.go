package routes

import (
	"encoding/json"
	"github.com/rutan/yarujan/lib/image"
	"github.com/rutan/yarujan/lib/uploader"
	"github.com/tuvistavie/securerandom"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

func Create(c web.C, w http.ResponseWriter, r *http.Request) {
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

	blob, err := createImage(contents)
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusBadRequest)
		return
	}

	client := uploader.New()
	key := generateKey() + ".jpg"
	url, err := client.UploadBlob(getFromEnv("AWS_S3_BUCKET_NAME", ""), key, blob)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(Image{Url: url})
}

func createImage(contents []byte) ([]byte, error) {
	canvas := image.LoadFromBlob(contents)
	defer canvas.Destroy()
	err := canvas.ResizeContain(480, 480)
	if err != nil {
		return nil, err
	}
	width, height := canvas.GetSize()
	text := selectText()
	setting := initTextSetting()
	setting.Size = width / float64(utf8.RuneCountInString(text))
	canvas.DrawText(text, width/2, height*0.92, &setting)
	return canvas.ExportBlob(), nil
}

func selectText() string {
	TEXT_LIST := []string{
		"LGTM",
		"いいね！",
		"よさそう",
		"みました",
		"やるじゃん",
	}
	return TEXT_LIST[rand.Intn(len(TEXT_LIST))]
}

func initTextSetting() image.TextSetting {
	setting := image.NewTextSetting()
	setting.Font = getFromEnv("FONT_NAME", "./fonts/toroman.ttf")
	setting.FillColor = getFromEnv("FILL_COLOR", "#ffffff")
	setting.BorderColor = getFromEnv("BORDER_COLOR", "#444444")
	setting.BorderWidth = 4
	return setting
}

func getFromEnv(envName string, defaultValue string) string {
	if len(os.Getenv(envName)) > 0 {
		return os.Getenv(envName)
	}
	return defaultValue
}

func generateKey() string {
	reverseTime := 32503648140 - time.Now().Unix()
	random, _ := securerandom.Hex(10)
	return strconv.FormatInt(reverseTime, 10) + "_" + random
}
