package image

import (
	"gopkg.in/gographics/imagick.v1/imagick"
)

type TextSetting struct {
	Font        string
	Size        float64
	Weight      uint
	Alignment   imagick.AlignType
	FillColor   string
	BorderColor string
	BorderWidth float64
}

func NewTextSetting() TextSetting {
	ts := TextSetting{}
	ts.Font = "Meiryo"
	ts.Size = 32
	ts.Weight = 500
	ts.Alignment = imagick.ALIGN_CENTER
	ts.FillColor = "#ffffff"
	ts.BorderColor = "#000000"
	ts.BorderWidth = 5
	return ts
}
