package avatar

import (
	"github.com/kolesa-team/go-webp/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func parseImage(contentType string, reader io.Reader) (img image.Image, err error) {
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "image/png":
		img, err = png.Decode(reader)
	case "image/webp":
		img, err = webp.Decode(reader, nil)
	case "image/gif":
		img, err = gif.Decode(reader)
	default:
		img, err = jpeg.Decode(reader)
	}

	return
}
