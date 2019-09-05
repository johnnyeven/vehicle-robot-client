package modules

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

// HLine draws a horizontal line
func HLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func VLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(img *image.RGBA, x1, y1, x2, y2, width int, col color.Color) {
	for i := 0; i < width; i++ {
		HLine(img, x1, y1+i, x2, col)
		HLine(img, x1, y2+i, x2, col)
		VLine(img, x1+i, y1, y2, col)
		VLine(img, x2+i, y1, y2, col)
	}
}

func ConvertImageToMat(image image.Image) (mat gocv.Mat, err error) {
	bounds := image.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()

	bytes := make([]byte, 0, x*y)
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := image.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8))
			bytes = append(bytes, byte(g>>8))
			bytes = append(bytes, byte(r>>8))
		}
	}

	mat, err = gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
	return
}
