package image

import (
	"gopkg.in/gographics/imagick.v1/imagick"
)

type Canvas struct {
	mw *imagick.MagickWand
}

func LoadFromBlob(blob []byte) Canvas {
	canvas := Canvas{}
	canvas.mw = createFreshImage(blob)
	canvas.mw.SetFormat("jpg")
	return canvas
}

func (self Canvas) Destroy() {
	self.mw.Destroy()
	self.mw = nil
}

func (self Canvas) GetSize() (float64, float64) {
	w := float64(self.mw.GetImageWidth())
	h := float64(self.mw.GetImageHeight())
	return w, h
}

func (self Canvas) ResizeContain(width float64, height float64) error {
	w, h := self.GetSize()
	ratioWidth := width / w
	ratioHeight := height / h

	var ratio float64
	if ratioWidth < ratioHeight {
		ratio = ratioWidth
	} else {
		ratio = ratioHeight
	}
	if ratio > 1 {
		ratio = 1
	}

	return self.mw.ResizeImage(uint(w*ratio), uint(h*ratio), imagick.FILTER_BOX, 1)
}

func (self Canvas) DrawText(text string, x float64, y float64, setting *TextSetting) {
	mw := self.mw
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	dw.SetFont(setting.Font)
	dw.SetFontSize(setting.Size)
	dw.SetTextAlignment(setting.Alignment)
	dw.SetFontWeight(setting.Weight)

	// border
	pw.SetColor(setting.BorderColor)
	dw.SetFillColor(pw)
	borderWidth := setting.BorderWidth
	mw.AnnotateImage(dw, x+borderWidth, y, 0, text)
	mw.AnnotateImage(dw, x-borderWidth, y, 0, text)
	mw.AnnotateImage(dw, x, y+borderWidth, 0, text)
	mw.AnnotateImage(dw, x, y-borderWidth, 0, text)
	mw.AnnotateImage(dw, x+borderWidth, y+borderWidth, 0, text)
	mw.AnnotateImage(dw, x-borderWidth, y+borderWidth, 0, text)
	mw.AnnotateImage(dw, x+borderWidth, y-borderWidth, 0, text)
	mw.AnnotateImage(dw, x-borderWidth, y-borderWidth, 0, text)

	// fill
	pw.SetColor(setting.FillColor)
	dw.SetFillColor(pw)
	mw.AnnotateImage(dw, x, y, 0, text)
}

func (self Canvas) ExportBlob() []byte {
	return self.mw.GetImageBlob()
}

func createFreshImage(blob []byte) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	tmw := imagick.NewMagickWand()
	defer tmw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()

	tmw.ReadImageBlob(blob)

	pw.SetColor("#ffffff")
	dw.SetFillColor(pw)
	mw.NewImage(tmw.GetImageWidth(), tmw.GetImageHeight(), pw)

	mw.CompositeImage(tmw, tmw.GetImageCompose(), 0, 0)
	return mw
}
